

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


/*
	`json:"command"`

		1. When unmarshaling JSON → map the JSON key "command" to this struct field.
		2. When marshaling struct to JSON → use "command" as the key.

	unmarshal -> JSON to Go structs
	marshal   -> Go structs to JSON


*/

// represents the struct how we will receive request for scheduling!!
type CommandRequest struct {
	Command     	string `json:"command"`
	ScheduledAt		string `json:"scheduled_at"`  //ISO8601 Format...
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

// 													---------------------- SERVER STARTER -------------------
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

	// return await shutdown
	return s.awaitShutdown()
}

// 					-------------------- SERVER SERVICES ------------------

// Handles if any new tasks needs to be added -> POST Req
func (s *SchedulerServer) handleScheduleTasks(w http.ResponseWriter, r *http.Request){

	// checks if the received req is post if not throws an err
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// now we need to prcess the req body

	var commandReq  CommandRequest

	if err := json.NewDecocder(r.Body).Decode(&commandReq); err != nil{
		// if we are not able to parse the json..raise its a bad request...
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if we are able to parse the JSON the request was valid..
	log.Printf("Received the Schedule Request: %+v", commandReq)

	// Parse the scheduled_at time
	// since we are expecting the time to be off ISO 8601 format..if the given format doesn't comply it raise an error!!
	scheduledTime, err := time.Parse(time.RFC3339, commandReq.ScheduledAt)
	if err != nil {
		http.Error(w, "Invalid date format. Use ISO 8601 format.", http.StatusBadRequest)
		return
	}

	/*
		scheduledTime.Unix() returns the number of seconds since January 1, 1970 UTC (a Unix timestamp).

		Input: "2025-08-22T10:30:00Z" (from API client)

		Parse: time.Parse(...) → time.Time

		Convert: Unix() → int64 seconds (e.g., 1755868200)

	*/

	// Convert the scheduled time to Unix timestamp
	unixTimestamp := time.Unix(scheduledTime.Unix(), 0)

	// once converted to desired timestamp -> We can now insert it into the DB :)


	taskId, err := s.insertTaskIntoDB(context.Background(), Task{Command: commandReq.Command, ScheduledAt: pgtype.Timestamp{Time:unixTimestamp}})

	// check for errors when inserting into the db
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to submit task. Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// create a response..we are creating an anon struct on the go and declaring value to it.
	response := struct {
    Command     string `json:"command"`
    ScheduledAt int64  `json:"scheduled_at"`
    TaskID      string `json:"task_id"`
	}{
    Command:     commandReq.Command,
    ScheduledAt: unixTimestamp.Unix(),
    TaskID:      taskId,
	}


	// since the res is in GO struct format -> Marshal it into JSON format
	jsonResponse, err := json.Marshal(response)

	// if err when marshalling
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)




}

// handle get task Status
func(s *SchedulerServer) handleGetTaskStatus(w http.ResponseWriter, r *http.Request){

	// http package is a gray area to me right now..once i learn about it..will wirte more desctiption about it.
	if r.Method != "GET"{
		http.Error(w, "Only Get Method is Allowed...", http.StatusMethodNotAllowed)
		return
	}

	// Get the taskid from query params...
	taskId := r.URL.Query().Get("task_id")

	// Check if the task ID is empty
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	// Query the DB to get the status
	var task Task

	err :=  s.dbPool.QueryRow(
			context.Background(),
			"SELECT * FROM tasks WHERE id = $1", taskId)
			.Scan(&task.Id, &task.Command, &task.ScheduledAt, &task.PickedAt, &task.StartedAt, &task.CompletedAt, &task.FailedAt)
	
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get task status. Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Prepare the Resp JSON
	/*	
		TaskId string `json:"task_id"` -> this line simply says when we are marshalling this TaskId should be put under task_id key in JSON
		omitempty -> If the field is empty (for a string, that means ""), then the field is omitted from the JSON output.
	*/

	response := struct{
		TaskId      string `json:"task_id"`,
		Command     string `json:"command"`
		ScheduledAt string `json:"scheduled_at, omitempty"`
		PickedAt    string `json:"picked_at,omitempty"`
		StartedAt   string `json:"started_at,omitempty"`
		CompletedAt string `json:"completed_at,omitempty"`
		FailedAt    string `json:"failed_at,omitempty"`
	}{
		TaskID:      task.Id,
		Command:     task.Command,
		ScheduledAt: "",
		PickedAt:    "",
		StartedAt:   "",
		CompletedAt: "",
		FailedAt:    "",
	}

	// Set the scheduled_at time if non-null.
	if task.ScheduledAt.Status == pgtype.Present {
		response.ScheduledAt = task.ScheduledAt.Time.String()
	}

	// Set the picked_at time if non-null.
	if task.PickedAt.Status ==  pgtype.Present {
		response.PickedAt = task.PickedAt.Time.String()
	}

	// Set the started_at time if non-null.
	if task.StartedAt.Status ==  pgtype.Present {
		response.StartedAt = task.StartedAt.Time.String()
	}

	// Set the completed_at time if non-null.
	if task.CompletedAt.Status ==  pgtype.Present {
		response.CompletedAt = task.CompletedAt.Time.String()
	}

	// Set the failed_at time if non-null.
	if task.FailedAt.Status ==  pgtype.Present {
		response.FailedAt = task.FailedAt.Time.String()
	}

	// Once the resp is ready mashell the struct to JSON

	// Convert the response struct to JSON

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	//set the Content-Type Header to JSON
	w.Header().Set("Content-Type","application/json")

	// write the json resp
	w.Write(jsonResponse)


}

// insert Tasks into the DB.....
// insertTaskIntoDB inserts a new task into the tasks table and returns the autogenerated UUID.
func(s *SchedulerServer) insertTaskIntoDB(ctx context.Context, task Task) (string, error)    {

	// sql statement
	sqlStatement := "INSERT INTO tasks (command, scheduled_at) VALUS ($1, $2) RETURNING id"

	var insertedId string

	err := s.dbPool.QueryRow(ctx, sqlStatement, task.Command, task.ScheduledAt.Time).Scan(&insertedId)

	if err != nil {
		return "", err
	}

	// Return the autogenerated UUID
	return insertedId, nil
}



//					----------------- SERVER STOPPER -----------------
// await shutdown method
// This is a method on the SchedulerServer type.
func (s *SchedulerServer ) awaitShutdown error {
	// create a buffered channel where process can check if they need to stop
	stop := make(chan os.Signal, 1)

	/*

		signal.Notify is a function from Go's os/signal package.

		Its job: register the channel (stop) to receive notifications when the program gets certain OS signals.

		After this call, whenever the OS sends one of the listed signals, Go will send that signal into the channel.

		syscall.SIGTERM:
			- Sent when the system or another process asks the program to terminate gracefully.

		os.Interrupt:
			- Usually triggered by Ctrl+C in the terminal.

	*/
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// blocks the code until somevalue is put into the server
	// since this is the receiver it will block..buffered channel is non blocking only for the sender
	<- stop

	return s.Stop()
}



func (s *SchedulerServer) Stop() error {

	// once the stop method is called -> we need to close the connection
	// first step is to close the dbPool
	s.dbPool.Close()

	// once db is closed we need to close the Http Server that was listening for requests
	// this is the standard GO way to gracefully shutdown the server!!
	if s.httpServer != nil{
		// gives 10 secs for all the remaining taks to be completed!!
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)

		defer cancel()

		return s.httpServer.Shutdown(ctx)

	}

	log.Println("Scheduler server and database pool stopped!!")

	return nil
}