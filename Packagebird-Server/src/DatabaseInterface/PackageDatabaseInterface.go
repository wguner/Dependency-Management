package DatabaseInterface

import (
	"context"
	"fmt"
	"log"
	structures "packagebird-server/src/structures"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

		// Increment version number of project
		if err := IncrementProjectVersion(client, newPackage.Name); err != nil {
			log.Printf("Error countered with incremented project version attached to package %v", newPackage.Name)
			return false, err
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

// Deletes the first document that matches the filter
func DeletePackage(client mongo.Client, newPackage structures.Package) (bool, error) {
	collection := client.Database("packagebird").Collection("packages")

	filter := bson.M{
		"$and": []bson.M{
			{"name": newPackage.Name},
			{"version": newPackage.Version},
		},
	}

	_, err := (collection.DeleteOne(context.TODO(), filter))

	if err != nil {
		log.Printf("Error encountered deleting a package in database")
		return false, err
	}
	return true, nil
}

// Deletes all matching documents
func DeleteAllPackages(client mongo.Client, newPackage structures.Package) (bool, error) {
	collection := client.Database("packagebird").Collection("packages")

	filter := bson.M{
		"$and": []bson.M{
			{"name": newPackage.Name},
			{"version": newPackage.Version},
		},
	}

	_, err := (collection.DeleteMany(context.TODO(), filter))

	if err != nil {
		log.Printf("Error encountered deleting a package in database")
		return false, err
	}

	//Return success without any error.
	return true, nil
}

func LookupPackage(client mongo.Client, newPackage structures.Package) (bool, error) {
	collection := client.Database("packagebird").Collection("packages")

	var result structures.Package
	filter := bson.M{
		"$and": []bson.M{
			{"name": newPackage.Name},
			{"version": newPackage.Version},
		},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		log.Printf("Document doesn't exist in database.")
		return false, nil
	} else if err != nil {
		// Find error encountered
		log.Printf("Error encountered searching for package in database")
		return false, err
	} else {
		log.Printf("Package with name %v is in the database", newPackage.Name)
		return true, nil
	}
}

func GetPackageDependencies(client mongo.Client, newPackage structures.Package) ([]string, error) {
	collection := client.Database("packagebird").Collection("packages")

	var result structures.Package
	var dependenciesList []string

	filter := bson.M{
		"$and": []bson.M{
			{"name": newPackage.Name},
			{"version": newPackage.Version},
		},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)

	if err == mongo.ErrNoDocuments {
		log.Printf("Document doesn't exist in database.")
		return nil, nil
	} else if err != nil {
		// Find error encountered
		log.Printf("Error encountered searching for package in database")
		return nil, err
	}

	// Iterate through dependencies, add item
	for _, dependency := range result.Dependencies {
		pstring := GetPackageNameVersion(dependency)
		pname := pstring[0]
		pversion := pstring[1]
		itemString := fmt.Sprintf("%v-%v", pname, pversion)
		dependenciesList = append(dependenciesList, itemString)
	}

	if err != nil {
		log.Printf("Error encountered searching for dependencies")
		return dependenciesList, err
	}

	// Successfully find all dependencies for package.
	return dependenciesList, nil
}

func GetPackage(client mongo.Client, name string, version int) (*structures.Package, error) {
	collection := client.Database("packagebird").Collection("packages")
	var result structures.Package
	filter := bson.M{
		"$and": []bson.M{
			{"name": name},
			{"version": version},
		},
	}

	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		log.Printf("Document doesn't exist in database.")
		return nil, err
	} else if err != nil {
		// Find error encountered
		log.Printf("Error encountered searching for package in database")
		return nil, err
	}

	return &result, nil
}

func GetPackageDependenciesRecurse(c mongo.Client, pname string, pversion int, l *[]string) []string {
	// Gets package from database
	p, err := GetPackage(c, pname, pversion)
	if err != nil {
		log.Printf("%v", err)
		return nil
	}

	// Iterate over dependencies
	for _, dep := range p.Dependencies {

		// Generate package string address
		pstring := GetPackageNameVersion(dep)
		pname := pstring[0]
		pversion, _ := strconv.Atoi(pstring[1])

		// If list doesn't contain reference to package, add and recurse
		if !contains(l, dep) {
			*l = append(*l, dep)
			ppackage, _ := GetPackage(c, pname, pversion)
			if ppackage.Dependencies != nil {
				GetPackageDependenciesRecurse(c, pname, pversion, l)
			}
		}
	}

	// Terminate branch when no more dependencies present
	return *l
}

// Checks if the string contains a value
func contains(list *[]string, value string) bool {
	for _, i := range *list {
		if i == value {
			return true
		}
	}
	return false
}

// Splits on the '-' character in a package string
// Example: 'oreo-v1' becomes '['oreo', '1']'
func GetPackageNameVersion(pstring string) []string {
	return strings.Split(pstring, "-v")
}

func GetPackages(client mongo.Client) ([]string, error) {
	collection := client.Database("packagebird").Collection("packages")
	option := options.Find()
	var result structures.Package
	returns := []string{}

	cursor, err := collection.Find(context.TODO(), bson.D{}, option)
	if err != nil {
		log.Printf("Failed to retrieve list of packages")
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode package")
			return []string{}, err
		}
		returns = append(returns, fmt.Sprintf("%v-%v", result.Name, result.Version))
	}

	return returns, nil
}

func GetProjects(client mongo.Client) ([]string, error) {
	collection := client.Database("packagebird").Collection("project")
	options := options.Find()
	var results []structures.Project
	var returns []string

	cursor, err := collection.Find(context.TODO(), bson.M{}, options)
	if err != nil {
		log.Printf("Failed to connect to database of projects")
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var project structures.Project
		err := cursor.Decode(&project)
		if err != nil {
			log.Printf("Cannot decode to project struct")
		}
		results = append(results, project)
	}
	for _, p := range results {
		log.Printf("Found:\t%v", p.Name)
		returns = append(returns, p.Name)
	}
	return returns, nil
}

func GetMembers(client mongo.Client) ([]string, error) {

	collection := client.Database("packagebird").Collection("members")
	option := options.Find()
	var result structures.Member
	returns := []string{}

	cursor, err := collection.Find(context.TODO(), bson.M{}, option)
	if err != nil {
		log.Printf("Failed to retrieve list of members")
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		if err := cursor.Decode(&result); err != nil {
			log.Printf("Failed to decode member")
			return []string{}, err
		}
		returns = append(returns, result.Name)
	}

	return returns, nil
}

func GetAllPackages(client mongo.Client, name string, version int64) ([]string, error) {
	// Check if package found
	found, err := LookupPackage(client, structures.Package{Name: name, Version: version})
	if err != nil {
		return nil, err
	}

	// If not found, return nothing
	if !found {
		return nil, err
	}

	// Get all package dependencies
	packages, err := GetPackageDependencies(client, structures.Package{Name: name, Version: version})
	if err != nil {
		return nil, err
	}
	packages = append(packages, fmt.Sprintf("%v-v%v", name, version))

	// Return list of package and it's dependencies
	return packages, nil
}

// Return list of packages back to client, make request for each file and return and extract to client as formatted
// Each package should have formatted list of references to dependencies
// Client should request file for each item in list, extract and format as specified
// Unfortunately means more request and network exchanges, compared to bundling all items into single file and transmitting whole
// Probably better to have atomic operations for relability
