package schema

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateSchema(client mongo.Client) {
	// Connects to or creates the database 'packagebird' in the MongoDB client
	db := client.Database("packagebird")

	// Creates collections and validation schema
	for collection, schema := range requiredCollectionsAndSchemas {
		db.CreateCollection(context.Background(), collection, options.CreateCollection().SetValidator(schema))
	}
}
