package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	collections "packagebird-server/mongo/enums"
	"packagebird-server/mongo/structures"
	"reflect"
)

// --- Utils ---

func GetDocumentFromCollectionNameByObjectId(client mongo.Client, collectionName string, objectId primitive.ObjectID) (*mongo.SingleResult, error) {
	collection := client.Database("packagebird").Collection(collectionName)
	filter := bson.M{
		"_id": objectId,
	}

	return collection.FindOne(context.Background(), filter), nil
}

func GetDocumentsFromCollectionName(client mongo.Client, collectionName string) (*mongo.Cursor, error) {
	collection := client.Database("packagebird").Collection(collectionName)
	result, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error retrieving many documents from collecton %v: %v", collectionName, err)
		return nil, err
	}
	defer result.Close(context.Background())

	return result, nil
}

func GetObjectFromCollectionNameByObjectId(client mongo.Client, collectionName string, objectId primitive.ObjectID, decodeType interface{}) (interface{}, error) {
	document, err := GetDocumentFromCollectionNameByObjectId(client, collectionName, objectId)
	if err != nil {
		return nil, err
	}
	result := decodeType
	err = document.Decode(result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func GetObjectsFromCollectionName(client mongo.Client, collectionName string, decodeType interface{}) (interface{}, error) {
	documents, err := GetDocumentsFromCollectionName(client, collectionName)
	if err != nil {
		return nil, err
	}
	results := reflect.New(reflect.SliceOf(reflect.TypeOf(decodeType)))
	for documents.Next(context.Background()) {
		result := &decodeType
		err := documents.Decode(result)
		if err != nil {
			return nil, err
		}
		results = reflect.Append(results.Elem(), reflect.ValueOf(result))
	}
	return results.Interface(), nil
}

// --- Package Get ---

func GetPackageByObjectId(client mongo.Client, objectId primitive.ObjectID) (*structures.Package, error) {
	result, err := GetObjectFromCollectionNameByObjectId(client, collections.Packages.String(), objectId, structures.Package{})
	if err != nil {
		return nil, err
	}
	decoded := result.(structures.Package)
	return &decoded, nil
}

func GetPackageByNameAndVersion(client mongo.Client, name string, version int64) (*structures.Package, error) {
	collection := client.Database("packagebird").Collection("packages")
	var filter = &bson.M{
		"$and": &[]bson.M{
			{"name": name},
			{"version": version},
		},
	}
	document := collection.FindOne(context.Background(), filter)

	result := &structures.Package{}
	err := document.Decode(&result)
	if err != nil {
		log.Printf("Error decoding document into package by name and version: %v", err)
		return nil, err
	}
	return result, nil
}

func GetPackagesByName(client mongo.Client, name string) ([]structures.Package, error) {
	collection := client.Database("packagebird").Collection("packages")
	var filter = &bson.M{
		"name": name,
	}

	document, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer document.Close(context.Background())

	var results []structures.Package

	for document.Next(context.Background()) {
		var result structures.Package
		err := document.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func GetPackages(client mongo.Client) ([]structures.Package, error) {
	documents, err := GetObjectsFromCollectionName(client, collections.Packages.String(), structures.Package{})
	if err != nil {
		return nil, err
	}
	return documents.([]structures.Package), nil
}

// --- Package Metadata Get ---

func GetPackageMetadataByNameAndVersion(client mongo.Client, name string, version int64) (*structures.PackageMetadata, error) {
	// TODO: To be filled
	return nil, nil
}

func GetPackagesMetadata(client mongo.Client) ([]structures.PackageMetadata, error) {
	documents, err := GetObjectsFromCollectionName(client, collections.PackagesMetadata.String(), structures.PackageMetadata{})
	if err != nil {
		return nil, err
	}
	return documents.([]structures.PackageMetadata), nil
}

// --- User Get ---

func GetUserByObjectId(client mongo.Client, objectId primitive.ObjectID) (*structures.User, error) {
	document, err := GetObjectFromCollectionNameByObjectId(client, collections.Users.String(), objectId, structures.User{})
	if err != nil {
		return nil, err
	}
	result := document.(structures.User)
	return &result, nil
}

func GetUserByName(client mongo.Client, name string) ([]structures.User, error) {
	// TODO: To be filled
	return nil, nil
}

func GetUsers(client mongo.Client) ([]structures.User, error) {
	documents, err := GetObjectsFromCollectionName(client, collections.Users.String(), structures.User{})
	if err != nil {
		return nil, err
	}
	return documents.([]structures.User), nil
}

// --- Authentication Get ---

func GetAuthenticationByUserObjectId(client mongo.Client, objectId primitive.ObjectID) (*structures.Authentication, error) {
	object, err := GetObjectFromCollectionNameByObjectId(client, collections.Authentications.String(), objectId, structures.Authentication{})
	if err != nil {
		return nil, err
	}
	return object.(*structures.Authentication), nil
}

func GetAuthentications(client mongo.Client) ([]structures.Authentication, error) {
	objects, err := GetObjectsFromCollectionName(client, collections.Authentications.String(), structures.Authentication{})
	if err != nil {
		return nil, err
	}
	return objects.([]structures.Authentication), nil
}

// --- Source Get ---

func GetSourceByObjectId(client mongo.Client, objectId primitive.ObjectID) (*structures.Source, error) {
	object, err := GetObjectFromCollectionNameByObjectId(client, collections.Sources.String(), objectId, structures.Source{})
	if err != nil {
		return nil, err
	}
	return object.(*structures.Source), nil
}

func GetSources(client mongo.Client) ([]structures.Source, error) {
	objects, err := GetObjectsFromCollectionName(client, collections.Sources.String(), structures.Source{})
	if err != nil {
		return nil, err
	}
	return objects.([]structures.Source), nil
}

// --- Project Get ---

func GetProjectByObjectId(client mongo.Client, objectId primitive.ObjectID) (*structures.Project, error) {
	object, err := GetObjectFromCollectionNameByObjectId(client, collections.Projects.String(), objectId, structures.Project{})
	if err != nil {
		return nil, err
	}
	return object.(*structures.Project), nil
}

func GetProjectByName(client mongo.Client, name string) (*structures.Project, error) {
	// TODO: To be filled
	return nil, nil
}

func GetProjects(client mongo.Client) ([]structures.Project, error) {
	objects, err := GetObjectsFromCollectionName(client, collections.Projects.String(), structures.Project{})
	if err != nil {
		return nil, err
	}
	return objects.([]structures.Project), nil
}

// --- Script Get ---

func GetScriptByObjectId(client mongo.Client, objectId primitive.ObjectID) (*structures.Script, error) {
	object, err := GetObjectFromCollectionNameByObjectId(client, collections.Scripts.String(), objectId, structures.Script{})
	if err != nil {
		return nil, err
	}
	return object.(*structures.Script), nil
}

func GetScriptByName(client mongo.Client) (*structures.Script, error) {
	// TODO: To be filled
	return nil, nil
}

// --- Graph Get ---

func GetGraphByObjectId(client mongo.Client, objectId primitive.ObjectID) (*structures.Graph, error) {
	object, err := GetObjectFromCollectionNameByObjectId(client, collections.Graphs.String(), objectId, structures.Graph{})
	if err != nil {
		return nil, err
	}
	return object.(*structures.Graph), nil
}

func GetGraphByName(client mongo.Client, name string) (*structures.Graph, error) {
	// TODO: To be filled
	return nil, nil
}

func GetGraphs(client mongo.Client) ([]structures.Graph, error) {
	objects, err := GetObjectsFromCollectionName(client, collections.Graphs.String(), structures.Graph{})
	if err != nil {
		return nil, err
	}
	return objects.([]structures.Graph), nil
}
