package DatabaseInterface

import (
	"fmt"
	"log"
	buildtest "packagebird-server/src/NetworkInterface/buildtest"
	"packagebird-server/src/structures"
	"testing"
)

func TestGetPackageDependenciesRecurse(t *testing.T) {
	list := []string{}

	fmt.Printf("List starts as: %v\n", list)

	expected := []string{"apple-v0", "blueberry-v1", "pineapple-v3", "pear-v74"}
	mongoDBClient, _ := MongoDBServerConnect("mongodb://localhost:27017")
	GetPackageDependenciesRecurse(*mongoDBClient, "oreo", 0, &list)

	fmt.Printf("List is: %v\n", list)
	fmt.Printf("Expected is: %v\n", expected)

	for _, dep := range expected {
		if !contains(&list, dep) {
			t.Fatalf(`List was: %v, wanted: %v`, list, expected)
		}
	}
}

func TestNewPackage(t *testing.T) {
	mongoDBClient, _ := MongoDBServerConnect("mongodb://localhost:27017")
	var testPackage = buildtest.PackageInfo{Name: "oreo", Version: 0}
	var testStruct = structures.Package{Name: "oreo", Version: 0}
	NewPackage(*mongoDBClient, testStruct)
	_, err := GetPackage(*structures.GlobalMongoClient, "oreo", 0)
	if err != nil {
		log.Printf("Could not add a package: %v", (testPackage.GetName()))
	}
}

func TestDeletePackage(t *testing.T) {
	mongoDBClient, _ := MongoDBServerConnect("mongodb://localhost:27017")
	var testPackage = buildtest.PackageInfo{Name: "oreo", Version: 0}
	var testStruct = structures.Package{Name: "oreo", Version: 0}
	NewPackage(*mongoDBClient, testStruct)
	_, err := GetPackage(*structures.GlobalMongoClient, "oreo", 0)
	if err != nil {
		log.Printf("Test package is succesfully added: %v", (testPackage.GetName()))
	}

	DeletePackage(*mongoDBClient, testStruct)
	if err == nil {
		log.Printf("Package: %v is deleted.", (testPackage.GetName()))
	}
}

