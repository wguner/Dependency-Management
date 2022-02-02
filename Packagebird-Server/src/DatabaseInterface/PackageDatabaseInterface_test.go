package DatabaseInterface

import (
	"fmt"
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
