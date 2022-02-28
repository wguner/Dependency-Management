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
	option := options.Find().SetProjection(bson.D{{"name", 1}, {"_id", 0}})
	collection := client.Database("packagebird").Collection("packages")

	type PackageResults struct {
		Names []string `bson:"name"`
	}

	var returns *PackageResults
	cursor, err := collection.Find(context.TODO(), bson.D{}, option)
	if err != nil {
		log.Printf("Failed to retrieve list of packages")
		return nil, err
	}

	cursor.All(context.TODO(), &returns)

	for _, result := range returns.Names {
		log.Printf("Listing package: %v\n", result)
	}

	return returns.Names, nil
}

func GetProjects(client mongo.Client) ([]string, error) {
	collection := client.Database("packagebird").Collection("projects")
	var results []structures.Project
	var returns []string

	cursor, err := collection.Find(context.TODO(), bson.D{{}}, nil)
	if err != nil {
		log.Printf("Failed to connect to database of projects")
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		var project *structures.Project
		err := cursor.Decode(project)
		if err != nil {
			log.Printf("Cannot decode to project struct")
		}
		results = append(results, *project)
	}
	for _, p := range results {
		log.Printf("Found:\t%v", p.Name)
		returns = append(returns, p.Name)
	}
	return returns, nil
}

func GetMembers(client mongo.Client) ([]string, error) {
	option := options.Find().SetProjection(bson.D{{"name", 1}, {"_id", 0}})
	collection := client.Database("packagebird").Collection("members")

	type MemberResults struct {
		Names []string `bson:"name"`
	}

	var returns *MemberResults
	cursor, err := collection.Find(context.TODO(), bson.D{}, option)
	if err != nil {
		log.Printf("Failed to retrieve list of members")
		return nil, err
	}

	cursor.All(context.TODO(), &returns)
	for _, result := range returns.Names {
		log.Printf("Listing Member: %v\n", result)
	}

	return returns.Names, nil
}
