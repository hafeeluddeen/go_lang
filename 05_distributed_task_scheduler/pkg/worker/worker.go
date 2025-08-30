package worker

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	taskProcessTime = time.Second * 5
	workerPoolSize  = 5 // no of workers in the pool
)

type WorkerServer struct {
	pb.UnimplementedWorkerServiceServer
	id         unint32
	serverPort string
	// address of the coordinator the worker needs to register itself with.
	coordinatorAddress string
	// listerner for the grpc server
	listener net.Listener
	//grpc server itself
	grpcServer *grpc.Server

	// obj to maintain GRPC connection
	coordinatorConnection *grpc.ClientConn

	coordinatorServiceClient pb.CoordinatorService

	heartBeatInterval time.Duration

	taskkQueue chan *pb.TaskRequest

	receivedTasks map[string]*pb.TaskRequest

	receivedTasksMutex sync.Mutex

	ctx    context.Context    // root context for all go routines
	cancel context.CancelFunc // cancel func to cancel all the go routines
	wg     sync.WaitGroup     // wg to wait for all the go routines to finish

}

// NewServer creates and returns a new WorkerServer.
func NewServer(port string, coordinator string) *WorkerServer {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerServer{
		id:                 uuid.New().ID(),
		serverPort:         port,
		coordinatorAddress: coordinator,
		heartbeatInterval:  common.DefaultHeartbeat,
		taskQueue:          make(chan *pb.TaskRequest, 100), // Buffered channel
		ReceivedTasks:      make(map[string]*pb.TaskRequest),
		ctx:                ctx,
		cancel:             cancel,
	}
}

func (w *WorkerServer) Start() error {

	// start the worker pool
	w.startWorkerPool(workerPoolSize)

	// connect with the co-ordinator
	if err := w.connectToCoordinator(); err != nil {
		return fmt.Errorf("failed to connect to coordinator: %w", err)
	}

	// close the GRPC connection
	defer w.closeGRPCConnection()

	// send periodic heartbeats...
	go w.periodicHeartbeat()

	// start the GRPC server...
	if err := w.StartGRPCServer(); err != nil {
		return fmt.Errorf("gRPC server failed....: %w", err)
	}

	return w.awaitShutdown()

}

func (w *WorkerServer) connectToCoordinator() error {
	log.Println("Connecting to Co-ordinator Server.....")

	var err error

	// dail the coordinator to connect with it...

	w.coordinatorConnection = grpc.Dial(w.coordinatorAddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())

	if err != nil {
		return err
	}

	// set the client from the coordinator to interact with it..
	w.coordinatorServiceClient = pb.NewCoordinatorServiceClient(w.coordinatorConnection)

	log.Println("Connected to coordinator!")
	return nil
}

func (w *WorkerServer) periodicHeartbeat() {
	// Add this goroutine to the waitgroup.
	w.wg.Add(1)

	// Signal this goroutine is done when the function returns
	defer w.wg.Done()

	// send heart beast at regular intervals of time..
	ticker := time.NewTicker(w.heartBeatInterval)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// once the time elapsed raise heartbeat
			if err := w.sendHeartbeat(); err != nil {
				log.Printf("Failed to send heartbeat: %v", err)
				return
			}
		// OR if the context is over shut down
		case <-w.ctx.Done():
			return
		}
	}

}

func (w *WorkerServer) sendHeartbeat() error {

	// get the address to send heatbeat...
	workerAddress := os.Getenv("WORKER_ADDRESS")

	if workerAddress == "" {
		// Fall back to using the listener address if WORKER_ADDRESS is not set
		workerAddress = w.listener.Addr().String()

	} else {
		workerAddress += w.serverPort
	}

	// once worker address has been found.. send the heart beat

	_, err := w.coordinatorServiceClient.SendHeartbeat(context.Background(), &pb.HeartbeatRequest{
		WorkerID: w.id,
		Address:  workerAddress,
	})

	return err

}

func (w *WorkerServer) StartGRPCServer() error {

	var err error

	if w.serverPort == "" {
		// Find a free port using a temporary socket
		w.listener, err = net.Listen("tcp", ":0")                                // Bind to any available port
		w.serverPort = fmt.Sprintf(":%d", w.listener.Addr().(*net.TCPAddr).Port) // Get the assigned port
	} else {
		// if the server port already exists listen into it.
		w.listener, err = net.Listen("tcp", w.serverPort)
	}

	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", w.serverPort, err)
	}

	log.Printf("Starting worker server on %s\n", w.serverPort)

	// start a new GRPC SERVER and register it..
	w.grpcServer = grpc.NewServer()
	pb.RegisterWorkerServiceServer(w.grpcServer, w)

	go func() {
		if err := w.grpcServer.Serve(w.listener); err != nil {
			log.Fatalf("GRPC Failed: %v", err)
		}
	}()

	return nil

}

// GRACE FULL SHUT DOWN....
func (w *WorkerServer) awaitShutdown() error {
	stop := make(chan os.Signal, 1)

	/*
		"Whenever the process receives SIGINT (Ctrl+C) or SIGTERM (typical kill command), send it into the stop channel."
	*/

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	return w.Stop()
}

// Stop Gracefully shuts down the Worker server...
func (w *WorkerServer) Stop() {

	// signals all the go routines to stop
	w.cancel()

	// waits for all the Go Routines to finish
	w.wg.Wait()

	w.closeGRPCConnection()
	log.Println("Worker server stopped")
	return nil
}

func (w *WorkerServer) closeGRPCConnection() {

	// close the grpc server
	if w.grpcServer != nil {
		w.grpcServer.GracefulStop()
	}

	// close the listener
	if w.listener != nil {
		if err := w.listener.Close(); err != nil {
			log.Printf("Error while closing the listener: %v", err)
		}
	}

	// close the coordinator connection
	if err := w.coordinatorConnection.Close(); err != nil {
		log.Printf("Error while closing client connection with coordinator: %v", err)

	}
}

// startWorkerPool starts a pool of worker goroutines.
func (w *WorkerServer) startWorkerPool(numWorkers int) {

	for i := 0; i < numWorkers; i++ {

		// add the worker to the wait group
		w.wg.Add(1)

		// instantiate the worker..
		go w.worker()
	}

}

// worker is the function executed by each worker GO Routine.
func (w *WorkerServer) worker() {

	// signal this worker is done when the function returns....
	defer w.wg.Done()

	for {
		select {
		// if a task is present int he Queue Take and run it..
		case task := <-w.taskQueue:
			go w.updateTaskStatus(task, pb.TaskStatus_STARTED)
			w.processTask(task)
			go w.updateTaskStatus(task, pb.TaskStatus_COMPLETE)

		case <-w.ctx.Done():
			return
		}
	}

}

// SubmitTask handles the submission of a task to the worker server.
func (w *WorkerServer) SubmitTask(ctx context.Context, req *pb.TaskRequest) (*pb.TaskResponse, error) {
	log.Printf("Received task: %+v", req)

	w.ReceivedTasksMutex.Lock()
	w.ReceivedTasks[req.GetTaskId()] = req
	w.ReceivedTasksMutex.Unlock()

	w.taskQueue <- req

	return &pb.TaskResponse{
		Message: "Task was submitted",
		Success: true,
		TaskId:  req.TaskId,
	}, nil
}

// updateTaskStatus updates the status of a task.
func (w *WorkerServer) updateTaskStatus(task *pb.TaskRequest, status pb.TaskStatus) {
	w.coordinatorServiceClient.UpdateTaskStatus(context.Background(), &pb.UpdateTaskStatusRequest{
		TaskId:      task.GetTaskId(),
		Status:      status,
		StartedAt:   time.Now().Unix(),
		CompletedAt: time.Now().Unix(),
	})
}

// processTask simulates task processing.
func (w *WorkerServer) processTask(task *pb.TaskRequest) {
	log.Printf("Processing task: %+v", task)
	time.Sleep(taskProcessTime)
	log.Printf("Completed task: %+v", task)
}
