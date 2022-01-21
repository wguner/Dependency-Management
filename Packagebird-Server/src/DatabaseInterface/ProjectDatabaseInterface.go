package DatabaseInterface

import (
	"context"
	"log"
	structures "packagebird-server/src/structures"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewProject(client mongo.Client, name string, description string) error {
	// First establish connection to database and collection
	collection := client.Database("packagebird").Collection("project")

	// Check that project isn't already present in the collection
	var result structures.Project
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "name", Value: name}}).Decode(&result)

	// !!! This conditional needs fixed, not elegant at all.
	if err != nil {
		if err == mongo.ErrNoDocuments {

			// If no project found with name, create new project
			log.Printf("Creating new project with name %v in database...", name)

			newProject := structures.Project{Name: name, Description: description, UUID: "", LatestVersion: 0, Packages: nil}
			_, insertErr := collection.InsertOne(context.TODO(), newProject)

			if insertErr != nil {
				log.Printf("Error encountered with inserting new project into database...")
				return insertErr
			}

			return nil
		} else {
			// Unknown error encountered
			log.Printf("Error encountered searching for project in database...")
			return err
		}
	} else {
		// Project already exist
		log.Printf("Project with name %v already exist in database...", name)
		return nil
	}
}
