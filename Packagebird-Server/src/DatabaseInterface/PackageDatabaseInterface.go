package DatabaseInterface

import (
	"context"
	"log"
	structures "packagebird-server/src/structures"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewPackage(client mongo.Client, newPackage structures.Package) (bool, error) {
	// First establish connection to database and collection
	collection := client.Database("packagebird").Collection("packages")

	// Check that duplicate package isn't already present
	var result structures.Package
	filter := bson.M{
		"$and": []bson.M{
			{"name": newPackage.Name},
			{"version": newPackage.Version},
		},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	// No duplicate package found, create a new entry
	if err == mongo.ErrNoDocuments {
		log.Printf("Creating new package with name %v in database", newPackage.Name)

		_, iErr := collection.InsertOne(context.TODO(), newPackage)
		if iErr != nil {
			log.Printf("Error encountered with inserting new package into database")
			return false, iErr
		}

		return true, nil
	} else if err != nil {
		// Find error encountered
		log.Printf("Error encountered searching for package in database")
		return false, err
	} else {
		// Package already exist
		log.Printf("Project with name %v and version %v already exist in database", newPackage.Name, newPackage.Version)
		return false, nil
	}
}
