package DatabaseInterface

import (
	"context"
	"fmt"
	"log"
	"os"
	configFile "packagebird-server/src/config"
	structures "packagebird-server/src/structures"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Create new project
func NewProject(client mongo.Client, name string, description string) (bool, error) {
	// First establish connection to database and collection
	collection := client.Database("packagebird").Collection("project")

	// Check that project isn't already present in the collection
	var result structures.Project
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "name", Value: name}}).Decode(&result)

	// !!! This conditional needs fixed, not elegant at all.
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Project source directory path
			projectSourcePath := configFile.Config.ProjectSourcePath
			nameSourcePath := fmt.Sprintf("%v/%v", projectSourcePath, name)
			if err := os.Mkdir(nameSourcePath, 0755); err != nil {
				log.Printf("Failed to create project directory at path specified\n%v", err)
				return false, err
			}

			// If no project found with name, create new project
			log.Printf("Creating new project with name %v in database...", name)

			newProject := structures.Project{Name: name, Description: description, UUID: "", LatestVersion: 0, Packages: nil, SourceFile: nameSourcePath}
			_, insertErr := collection.InsertOne(context.TODO(), newProject)

			if insertErr != nil {
				log.Printf("Error encountered with inserting new project into database...")
				return false, insertErr
			}

			return true, nil
		} else {
			// Unknown error encountered
			log.Printf("Error encountered searching for project in database...")
			return false, err
		}
	} else {
		// Project already exist
		log.Printf("Project with name %v already exist in database...", name)
		return false, nil
	}
}

// Lookup particular project
func LookupProject(client mongo.Client, name string, description string) (bool, error) {
	collection := client.Database("packagebird").Collection("project")

	log.Printf("Looking up project with name: %v", name)
	_, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "name", Value: name}})

	if err == mongo.ErrNoDocuments {
		log.Printf("Project with name %v is not in the database", name)
		return false, nil
	} else if err == nil {
		return true, nil
	} else {
		// Find error encountered
		log.Printf("Error encountered searching for project in database")
		return false, err
	}
}

// Update project attributes
func UpdateProject(c mongo.Client, name string, dependencies []string) error {
	collection := c.Database("packagebird").Collection("project")

	var result structures.Project
	filter := bson.M{"name": name}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	result.Packages = append(result.Packages, dependencies...)

	_, err = collection.UpdateOne(context.TODO(), bson.M{"name": name}, bson.D{{"$set", bson.M{"dependencies": result.Packages}}})
	if err != nil {
		log.Printf("%v", err)
		return err
	} else {
		log.Printf("Updated dependencies for project: %v", name)
		return nil
	}
}

// Increment version of project
func IncrementProjectVersion(c mongo.Client, name string) error {
	collection := c.Database("packagebird").Collection("project")
	project, err := GetProject(*collection, name)
	if err != nil {
		log.Printf("Cannot find and retrieve project %v from MongoDB", name)
		return err
	}

	version := project.LatestVersion + 1
	_, err = collection.UpdateOne(context.TODO(), bson.M{"name": name}, bson.D{{"$set", bson.M{"version": version}}})

	if err != nil {
		log.Printf("Error incrementing project %v version: %v", name, err)
		return err
	}
	return nil
}

// Get particular project
func GetProject(c mongo.Collection, name string) (*structures.Project, error) {
	var result structures.Project
	filter := bson.M{"name": name}
	if err := c.FindOne(context.TODO(), filter).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Get member from database
func GetMember(client mongo.Client, name string) (*structures.Member, error) {
	collection := client.Database("packagebird").Collection("members")
	var member structures.Member
	filter := bson.M{"name": name}
	if err := collection.FindOne(context.TODO(), filter).Decode(&member); err != nil {
		return nil, err
	}
	return &member, nil
}

// Create member in database
func CreateMember(client mongo.Client, member structures.Member) error {
	collection := client.Database("packagebird").Collection("members")
	// If member not already in database, proceed operation
	_, err := GetMember(client, member.Name)
	if err == mongo.ErrNoDocuments {
		_, err := collection.InsertOne(context.TODO(), member)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get Admin user from database
func GetAdminMember(client mongo.Client) (*structures.Member, error) {
	// Find admin user from database
	admin, err := GetMember(client, "admin")
	if err == mongo.ErrNoDocuments {
		// Create administrator member
		admin := &structures.Member{
			Name:     "admin",
			Password: "password",
			Level:    "ADMIN",
			Employed: true,
		}
		if err := CreateMember(client, *admin); err != nil {
			return nil, err
		}
		return admin, nil
	} else if err != nil {
		// Return error
		return nil, err
	} else {
		// Return adminstrator member
		return admin, nil
	}
}

// Remove member in database
func RemoveMember(client mongo.Client, name string) error {
	collection := client.Database("packagebird").Collection("members")
	// If member not already in database, cancel operation
	_, err := GetMember(client, name)
	if err == mongo.ErrNoDocuments {
		log.Printf("Attempting to delete member not in database with name: %v", name)
		return err
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"name": name})
	if err != nil {
		return err
	}

	log.Printf("Successfully removed member %v from server database.", name)
	return nil
}

// Remove dependencies from a project
func RemoveDependency(client mongo.Client, pname string, depname string, depversion int64) error {
	// Remove from a particular project
	collection := client.Database("packagebird").Collection("project")

	// If project not in database, cancel operation, else map to project
	var project structures.Project
	if err := collection.FindOne(context.TODO(), bson.M{"name": pname}).Decode(&project); err != nil {
		return nil
	}

	// Project created and loaded, if dependency in project, remove, else signal and proceed
	dependencies := project.Packages
	depstring := fmt.Sprintf("%v-v%v", depname, depversion)
	for i, dep := range dependencies {
		if dep == depstring {
			dependencies = append(dependencies[:i], dependencies[i+1:]...)
		}
	}

	// Update project dependencies list
	project.Packages = dependencies

	// Update database reference
	_, err := collection.ReplaceOne(context.Background(), bson.M{"name": pname}, project)
	if err != nil {
		return err
	}
	return nil
}
