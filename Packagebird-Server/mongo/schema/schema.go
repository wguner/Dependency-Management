package schema

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateSchema(client mongo.Client) error {
	// Connects to or creates the database 'packagebird' in the MongoDB client
	db := client.Database("packagebird")

	// Creates collections and validation schema
	for collection, schema := range requiredCollectionsAndSchemas {
		err := db.CreateCollection(context.Background(), collection, options.CreateCollection().SetValidator(schema))
		if err != nil && !errors.As(err, &mongo.CommandError{}) {
			return err
		}
	}
	return nil
}
