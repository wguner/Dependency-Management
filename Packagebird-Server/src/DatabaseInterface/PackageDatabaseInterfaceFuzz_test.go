package DatabaseInterface

import (
	"log"
	buildtest "packagebird-server/src/NetworkInterface/buildtest"
	"packagebird-server/src/structures"
	"testing"
)

func FuzzNewPackage(f *testing.F) {
	mongoDBClient, _ := MongoDBServerConnect("mongodb://localhost:27017")
	var testPackage = buildtest.PackageInfo{Name: "oreo", Version: 0}
	var testStruct = structures.Package{Name: "oreo", Version: 0}
	NewPackage(*mongoDBClient, testStruct)
	_, err := GetPackage(*structures.GlobalMongoClient, "oreo", 0)
	if err != nil {
		log.Printf("Could not add a package: %v", (testPackage.GetName()))
	}
}

func FuzzDeletePackage(f *testing.F) {
	mongoDBClient, _ := MongoDBServerConnect("mongodb://localhost:27017")
	var testPackage = buildtest.PackageInfo{Name: "oreo", Version: 0}
	var testStruct = structures.Package{Name: "oreo", Version: 0}
	NewPackage(*mongoDBClient, testStruct)
	_, err := GetPackage(*structures.GlobalMongoClient, "oreo", 0)
	if err != nil {
		log.Printf("Could not add a package: %v", (testPackage.GetName()))
	}
	DeletePackage(*mongoDBClient, testStruct)
	if err == nil {
		log.Printf("Package: %v is deleted.", (testPackage.GetName()))
	}
}

func FuzzDeleteAllPackages(f *testing.F) {
	mongoDBClient, _ := MongoDBServerConnect("mongodb://localhost:27017")
	var testPackage = buildtest.PackageInfo{Name: "oreo", Version: 0}
	var testStruct = structures.Package{Name: "oreo", Version: 0}
	NewPackage(*mongoDBClient, testStruct)
	_, err := GetPackage(*structures.GlobalMongoClient, "oreo", 0)
	if err != nil {
		log.Printf("Could not add a package: %v", (testPackage.GetName()))
	}
	DeleteAllPackages(*mongoDBClient, testStruct)
	if err == nil {
		log.Printf("Package: %v is deleted.", (testPackage.GetName()))
	}
}
