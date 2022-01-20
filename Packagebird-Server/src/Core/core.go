package main

import (
	"log"
	NetworkInterface "packagebird-server/src/networkinterface"
)

// Temporary constant for address
const ADDRESS = "127.0.0.1:50051"

// Entry-point for the server application
func main() {
	if err := NetworkInterface.PackagebirdServerStart(ADDRESS); err != nil {
		log.Fatalf("Error encountered on gRPC server:\n%v\nShutting down...", err)
	}
}
