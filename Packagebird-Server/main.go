package main

import (
	"fmt"
	"log"
	"packagebird-server/const"
	"packagebird-server/filesystem"
	"packagebird-server/global"
	"packagebird-server/mongo/schema"
	"packagebird-server/src/DatabaseInterface"
	"packagebird-server/src/NetworkInterface"
)

// Entry-point for packagebird server
func main() {
	var err error

	// Connect to MongoDB service
	global.GlobalMongoClient, err = DatabaseInterface.MongoDBServerConnect(_const.DefaultMongoUri)
	if err != nil {
		log.Fatal(err)
	}

	// Create default subdirectories if not already created
	err = filesystem.CreateDefaultSubdirectories()
	if err != nil {
		log.Fatal(err)
	}

	// Create default schema, will display errors if already created
	err = schema.CreateSchema(*global.GlobalMongoClient)
	if err != nil {
		log.Fatal(err)
	}

	// Begin listening for incoming request
	err = NetworkInterface.PackagebirdServerStart(fmt.Sprintf("%v:%v", _const.DefaultPublicListenAddress, _const.DefaultListenAddressPort), global.GlobalMongoClient)
	if err != nil {
		log.Fatal(err)
	}

	// Disconnect from MongoDB service
	err = DatabaseInterface.MongoDBServerDisconnect(*global.GlobalMongoClient)
	if err != nil {
		log.Fatal(err)
	}
}

// Server defaults to storing package sources in ./packages
// Packages are stored under ./packages/{version}/{name}
// Server defaults to storing project sources in ./projects
// Projects are stored under ./projects/{name}
// Server defaults to creating single admin user if not already created with 'admin' 'password'
