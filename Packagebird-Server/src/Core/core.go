package main

import (
	"log"
	"packagebird-server/src/DatabaseInterface"
	NetworkInterface "packagebird-server/src/networkinterface"
)

// Temporary constant for address
const (
	ADDRESS  = "127.0.0.1:50051"
	MONGOURI = "mongodb//:localhost:27017"
)

// Entry-poin for the server application
func main() {
	mongoDBClient, err := DatabaseInterface.MongoDBServerConnect(MONGOURI)
	if err != nil {
		log.Fatalf("Error encountered on server connecting to MongoDB:\n%v\nShutting down...", err)
	}

	if err := NetworkInterface.PackagebirdServerStart(ADDRESS); err != nil {
		log.Fatalf("Error encountered on gRPC server:\n%v\nShutting down...", err)
	}

	err = DatabaseInterface.MongoDBServerDisconnect(mongoDBClient);
	if err != nil {
		log.Fatalf("Error encountered on server disconnecting from MongoDB:\n%v\nShutting down...", err)
	}
}
