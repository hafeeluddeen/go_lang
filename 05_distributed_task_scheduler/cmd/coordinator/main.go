package main

import "flag"

const (
	coordinatorPort = flag.String("coordinator_port", ":8080", "Port on which the Coordinator servers requests....")
)

func main() {

	flag.Parse()
	dbConnectionString := common.GetDBConnectionString()

	coordinator := coordinator.NewServer(*coordinatorPort, dbConnectionString)

	coordinator.Start()

}
