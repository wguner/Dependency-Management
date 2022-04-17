package DatabaseInterface

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Create and return a MongoDB client connection to a database specified by a URI
func MongoDBServerConnect(uri string) (*mongo.Client, error) {

	// Makes a 'client' connection to the MongoDB server service, passed as the string URI. Logs fatal error if process somehow incomplete.
	mongodbClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Printf("Error encountered with initial connection to %v...", uri)
		return nil, err
	}

	// Conducts an initial ping operation to check if connection to MongoDB server service successful.
	err = mongodbClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Printf("ERror encountered with initial ping to %v...", uri)
	}

	// Assuming above steps are successful, notify server terminal.
	fmt.Println("Successfully connected to local MongoDB server...")

	return mongodbClient, nil
}

// Disconnect a passed MongoDB client from a MongoDB connection
func MongoDBServerDisconnect(client mongo.Client) error {
	err := client.Disconnect(context.Background())
	if err != nil {
		log.Printf("Error encountered with disconnection from MongoDB server service...")
		return err
	}

	log.Printf("Successfully disconnected from MongoDB server...")
	return nil
}
