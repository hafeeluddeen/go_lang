

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
}