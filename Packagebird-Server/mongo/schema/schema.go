package schema

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Schema(client mongo.Client) {
	// Connects to or creates the database 'packagebird' in the passed MongoDB client
	db := client.Database("packagebird")

	// Creates the 'packages' collection with validation
	db.CreateCollection(context.Background(), "packages", options.CreateCollection().SetValidator(packageSchema))

}
