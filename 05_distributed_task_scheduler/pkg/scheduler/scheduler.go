

// Scheduler will receive command
type CommandRequest struct {
	Command     string `json:"command"`
	ScheduledAt string `json:"command"`
}

// Tasks Structure
type Task struct {
	Id          string
	Command     string
	ScheduledAt pgtype.Timestamp
	PickedAt    pgtype.Timestamp
	StartedAt   pgtype.Timestamp
	CompletedAt pgtype.Timestamp
	FailedAt    pgtype.Timestamp
}

// Scheduler Server -> An HTTP server that manages tasks
type SchedulerServer struct {
	serverPort         string
	dbConnectionString string
	dbPool             *pgxpool.Pool
	ctx                context.Context
	cancel             context.CancelFunc
	httpServer         *http.Server
}

// create a new server and returns it/
func NewServer(port string, dbConnectionString string){

	// inorder to create the sub process created by this server we need to create a context

	ctx, cancel := context.WithCancel(context.Background())

	return &SchedulerServer{
		serverPort: port,
		dbConnectionString: dbConnectionString,
		ctx: ctx,
		cancel: cancel
	}
} 


// Start initilaizes and starts the server
func (s *SchedulerServer) Start() error{

	var err error

	s.dbPool, err := common.ConnectToDatabase(s.ctx, s.dbConnectionString)

	if err != nil{
		return err
	}

	// grey are for me right now
	http.HandleFunc("/schedule", s.handleScheduleTask)
	http.HandleFunc("/status/", s.handleGetTaskStatus )

	// grey are for me
	s.httpServer = &http.Server{
		Addr: s.serverPort,
	}

	log.PrintF("Starting the Scheduler Server on Port: %v\n",s.serverPort)


	/// start the scheduler server in a seperate goroutine......
	// ErrServerClosed means server was closed gracefully meaning it was stopped by us not due to some error.
	go func(){
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error %s\n", err)
		}
		
	}()
}