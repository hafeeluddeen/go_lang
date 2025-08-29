package coordinator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
	GRPC is low latency and high throughput when compared to regular

	http requests and uses protobuffs for serialization where as regular http uses json.

	So we are using grpc for for master and workers
*/

const (
	shutdownTimeout  = time.Second * 5
	defaultMaxMisses = 1
	scanInterval     = time.Second * 10
)

type CoordinatorServer struct {
	pb.UnimplementedCoordinatorServiceServer
	serverPort          string // port in which the server will run
	listener            net.Listener
	grpcServer          *grpc.Server
	WorkerPool          map[uint32]*workerInfo // all the workers associated with the co-ordinator
	WorkerPoolMutex     sync.Mutex             // locks the workerpool whenever we access it.
	WorkerPoolKeys      []uint32
	WorkerPoolKeysMutex sync.RWMutex
	maxHeartbeatMisses  uint8         // max heartbeat misses allowed for a worker before it gets de-registered
	heartbeatInterval   time.Duration // intervel at which each worker should send its heartbeat.
	roundRobinIndex     uint32        // helper var used to decide which worker is next in the round robin
	dbConnectionString  string
	dbPool              *pgxpool.Pool
	ctx                 context.Context    // The root context for all goroutines
	cancel              context.CancelFunc // Function to cancel the context
	wg                  sync.WaitGroup     // WaitGroup to wait for all goroutines to finish
}

type workerInfo struct {
	heartbeatMisses     uint8  // heartbeat misses a worker had
	address             string // network address of a worker
	grpcConnection      *grpc.ClientConn
	workerServiceClient pb.WorkerServiceClient // client to talk with this server
}

// creates a new co-ordinator server instances and returns it.
func NewServer(port string, dbConnectionString string) *CoordinatorServer {
	ctx, cancel := context.WithCancel(context.Background())

	return &CoordinatorServer{
		WorkerPool:         make(map[uint32]*workerInfo),
		maxHeartbeatMisses: defaultMaxMisses,
		heartbeatInterval:  common.DefaultHeartbeat,
		dbConnectionString: dbConnectionString,
		serverPort:         port,
		ctx:                ctx,
		cancel:             cancel,
	}
}

// Start inititaes the servers's operations

func (s *CoordinatorServer) Start() {
	var err error

	// this manages the workers..
	go s.manageWorkerPool()

	// start the GRPC server
	if err := s.StartGRPCServer(); err != nil {
		return fmt.Errorf("gRPC server start failed: %w", err)
	}

	// connect the GRPC server to the DB

	s.dbPool, err = common.ConnectToDatabase(s.ctx, s.dbConnectionString)

	if err != nil {
		return err
	}

	// this go routine scans the DB for any tasks that needs to be implemented in next 30 secs.
	// It fires queries continousy to check the tasks
	go s.scanDatabase()

	// blocks until we shutdown the coordinator server
	return s.awaitShutdown()
}

// HEARTBEAT REALTED FUNCTIONS --> STARTS
func (s *CoordinatorServer) SendHeartbeat(ctx context.Context, in *pb.HeartbeatRequest) *pb.HeartbeatResponse {
	/*
		when a worker sends heartbeat its either of 2 things.

		1. This is the first time the worker is sending the heartbeat we need to register it.
		2. If already registered we need to increment the no of times we have recorded the HeartBeat for that worker

	*/

	// to make sure n other chnages happen we accquire a lock on the workerpool mutex
	s.WorkerPoolMutex.Lock()

	// once all done release the lock
	defer s.WorkerPoolMutex.Unlock()

	// get the workerid from heartbeat req
	workerID := in.GetWorkerId()

	// if worker already exists..reset its heartbeat misses to 0
	if worker, ok := s.WorkerPool[workerID]; ok {
		// chnaging here will diretly affect it because we gave pointer to it -> WorkerPool map[uint32]*workerInfo
		worker.heartbeatMisses = 0
	} else {
		log.Println("Registering worker:", workerID)
		// register the worker in our pool
		conn, err := grpc.Dial(in.GetAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		defer conn.Close()
		if err != nil {
			return nil, err
		}

		s.WorkerPool[workerID] = &workerInfo{
			address:             in.GetAddress(),
			grpcConnection:      conn,
			workerServiceClient: pb.NewWorkerServiceClient(conn),
		}

		/*
			get a mutex lock on keys ->
				this keys help us to decide which worker next task needs to be assigned to in a round robin scheduling..

			We can't do it directly in the worker pool since keys in the map don't have a definite order

		*/
		s.WorkerPoolKeysMutex.Lock()
		defer s.WorkerPoolKeysMutex.Unlock()

		workerCount := len(s.WorkerPool)
		s.WorkerPoolKeys = make([]uint32, 0, workerCount)

		for k := range s.WorkerPool {
			s.WorkerPoolKeys = append(s.WorkerPoolKeys, k)
		}

		log.Println("Registered Worker --> ", workerID)
	}

	return &pb.HeartbeatResponse{Acknowledged: true}, nil

}

func (s *CoordinatorServer) manageWorkerPool() {

	/*
		Adding an waitgroup allows us to gracefully wait
		for the GO Routine to end when we are shutting down the server...
	*/

	s.wg.Add(1)
	defer s.wg.Done()

	/*
		This function runs in the background, and every few seconds,
		it removes workers who haven’t sent a heartbeat,
		and it stops when the program is shutting down.
	*/

	ticker := time.NewTicker(time.Duration(s.maxHeartbeatMisses) * s.heartbeatInterval)
	defer ticker.Stop()

	/*
		This loop constatly listens for any signal from a particular worker or from coordinator
		if it wants to remove a worker or stop the program completely if master wishes too
	*/
	for {
		select {
		// checks at regular intervals for inactive workers and removes them..
		case <-ticker.C:
			s.removeInactiveWorkers()
		// we receive done signal when we shut down the server...
		case <-s.ctx.Done():
			return
		}
	}

}

func (s *CoordinatorServer) removeInactiveWorkers() {
	/*
		We accquire lock on the worker pool because we might be making some changes on the worker pool
		[ removing few workers if they are inactive ]

	*/

	s.WorkerPoolMutex.Lock()
	defer s.WorkerPoolMutex.Unlock()

	for workerId, worker := range s.WorkerPool {
		if worker.heartbeatMisses > s.maxHeartbeatMisses {
			log.Printf("Removing Inactive worker....: %d\n", workerId)
			worker.grpcConnection.Close()
			delete(s.WorkerPool, workerId)

			/*
				accquire lock on woker pool keys mutex
				because we are going to make changes to the worker pool keys cause we had removed a worker
			*/

			s.WorkerPoolKeysMutex.Lock()

			workerCount := len(s.WorkerPool)
			s.WorkerPoolKeys = make([]uint32, 0, workerCount)

			for k := range s.WorkerPool {
				s.WorkerPoolKeys = append(s.WorkerPoolKeys, k)
			}

			s.WorkerPoolKeysMutex.Unlock()
		} else {
			worker.heartbeatMisses++
		}
	}

}

// HEARTBEAT REALTED FUNCTIONS --> ENDS

// TASK EXECUTION RELATED FUNCTIONS --> STARTS

func (s *CoordinatorServer) scanDB() {
	ticker := time.Ticker(scanInterval)

	// scanning the DB at regular intervals of time so see if we can get any tasks that needs to be executed
	for {
		select {
		case <-ticker.C:
			go executeAllScheduledTasks()

		case <-s.ctx.Done():
			log.Println("Shutting down database scanner.")
			return
		}
	}
}

func (s *CoordinatorServer) executeAllScheduledTasks() {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		log.Printf("Unable to start transaction: %v\n", err)
		return
	}

	// if there is an error rollback
	defer func() {
		/*
			Sometimes you might call Rollback after the transaction is already finished (committed or rolled back).

			In that case, Rollback often returns an error like "tx is closed".

			That error isn’t really a failure — it just means the transaction is already done.

			So the code ignores that specific error, but logs/handles any other rollback errors.
		*/
		if err := tx.Rollback(ctx); err != nil && err.Error() != "tx is closed" {
			log.Printf("ERROR: %#v", err)
			log.Printf("Failed to rollback transaction: %v\n", err)
		}
	}()

	rows, err := tx.Query(ctx, `SELECT id, command FROM tasks WHERE scheduled_at < (NOW() + INTERVAL '30 seconds') AND picked_at IS NULL ORDER BY scheduled_at FOR UPDATE SKIP LOCKED`)
	if err != nil {
		log.Printf("Error executing query: %v\n", err)
		return
	}
	defer rows.Close()

	// after getting all the tasks that needs to be scheduled...we put those taks into the TaskRequest Array..
	var tasks []*pb.TaskRequest

	for rows.Next() {
		var id, command string

		// rows.Scan copies the values from the current row into your variables.
		if err := rows.Scan(&id, &command); err != nil {
			log.Printf("Failed to scan row:%v\n", err)
			continue
		}

		tasks = append(tasks, &pb.TaskRequest{TaksId: id, Data: command})
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating rows: %v\n", err)
		return
	}

	// send the tasks to the workers for execution

	for _, task := range tasks {
		if err := s.submitTaskToWorker(task); err != nil {
			log.Printf("Failed to submit task %s: %v\n", task.GetTaskId(), err)
			continue
		}

		// once the task has been submitted update its picked_at time
		if _, err := tx.Exec(ctx, `UPDATE tasks SET picked_at = NOW() WHERE id = $1`, task.GetTaskId()); err != nil {
			log.Printf("Failed to update task %s: %v\n", task.GetTaskId(), err)
			continue
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Printf("Failed to commit transaction: %v\n", err)
	}
}

func (s *CoordinatorServer) getNextWorker() *workerInfo {

	// accquire read lock on the worker pool keys
	s.WorkerPoolKeysMutex.RLock()

	defer s.WorkerPoolKeysMutex.RUnlock()

	workerCount := len(s.WorkerPoolKeys)

	if workerCount == 0 {
		return nil
	}

	// get the worker via round robin method....
	worker := s.WorkerPool[s.WorkerPoolKeys[s.roundRobinIndex%uint32(workerCount)]]
	s.roundRobinIndex++
	return worker
}

func (s *CoordinatorServer) submitTaskToWorker(task *pb.TaskRequest) error {

	// get the next worker in line
	worker := s.getNextWorker()

	if worker == nil {
		return errors.New("no workers availbale...")
	}

	_, err := worker.workerServiceClient.SubmitTask(context.Background(), task)

	return err
}

// TASK EXECUTION RELATED FUNCTIONS --> ENDS
