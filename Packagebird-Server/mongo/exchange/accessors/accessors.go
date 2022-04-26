package accessors

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	collections "packagebird-server/mongo/enums"
	"packagebird-server/mongo/structures"
	"reflect"
)

// --- Utils ---

func ErrorCheck(err error) error {
	if err != nil {
		return err
	}
	return nil
}

func Filter[T any](arr []T, filter func(elem T) bool) []T {
	var r []T
	for i := range arr {
		if filter(arr[i]) {
			r = append(r, arr[i])
		}
	}
	return r
}

// --- Filter Utils ---

func filterByObjectId(objectId primitive.ObjectID) *bson.M {
	var filter = &bson.M{
		"_id": objectId,
	}
	return filter
}

func filterByNameAndVersion(name string, version int64) *bson.M {
	filter := &bson.M{
		"$and": []bson.M{
			bson.M{
				"name": name,
			},
			bson.M{
				"version": version,
			},
		},
	}
	return filter
}

func filterByName(name string) *bson.M {
	filter := &bson.M{
		"name": name,
	}
	return filter
}

// --- Get Utils ---

func GetDocumentFromCollectionNameByObjectId(client mongo.Client, collectionName string, objectId primitive.ObjectID) (*mongo.SingleResult, error) {
	collection := client.Database("packagebird").Collection(collectionName)
	filter := bson.M{
		"_id": objectId,
	}

	return collection.FindOne(context.Background(), filter), nil
}

func GetDocumentsFromCollectionName(client mongo.Client, collectionName string) (*mongo.Cursor, error) {
	collection := client.Database("packagebird").Collection(collectionName)
	result, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Printf("Error retrieving many documents from collecton %v: %v", collectionName, err)
		return nil, err
	}

	return result, nil
}

func GetObjectFromCollectionNameByObjectId(client mongo.Client, collectionName string, objectId primitive.ObjectID, decodeType interface{}) (interface{}, error) {
	document, err := GetDocumentFromCollectionNameByObjectId(client, collectionName, objectId)
	if err != nil {
		return nil, err
	}
	result := &decodeType
	err = document.Decode(*result)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

func GetObjectsFromCollectionName(client mongo.Client, collectionName string, decodeType interface{}) (interface{}, error) {
	documents, err := GetDocumentsFromCollectionName(client, collectionName)
	if err != nil {
		return nil, err
	}

	values := reflect.ValueOf(decodeType).Elem()
	slices := values.Slice(0, values.Cap())
	for documents.Next(context.Background()) {
		element := reflect.New(slices.Type().Elem())
		if err := documents.Decode(element.Interface()); err != nil {
			return nil, err
		}
		slices = reflect.Append(slices, element.Elem())
	}
	values.Set(slices)
	return values.Interface(), nil
}

func GetObjectsFromCollectionNameAndFilter(client mongo.Client, collectionName string, decodeType interface{}, filter interface{}) (interface{}, error) {
	collection := client.Database("packagebird").Collection(collectionName)
	documents, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer documents.Close(context.Background())

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

func GetObjectFromCollectionNameAndFilter(client mongo.Client, collectionName string, decodeType interface{}, filter interface{}) (interface{}, error) {
	collection := client.Database("packagebird").Collection(collectionName)
	document := collection.FindOne(context.Background(), filter)

	result := &decodeType
	err := document.Decode(*result)
	if err != nil {
		return nil, err
	}
	return *result, nil
}

// --- Set Utils ---

func SetDocumentInCollectionNameByObjectId(client mongo.Client, collectionName string, objectId primitive.ObjectID, document []byte) error {
	collection := client.Database("packagebird").Collection(collectionName)
	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": objectId}, document, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}
	return nil
}

func SetObjectInCollectionNameByObjectId(client mongo.Client, collectionName string, objectId primitive.ObjectID, toObject interface{}) error {
	marshall, err := bson.Marshal(toObject)
	if err != nil {
		return err
	}
	err = SetDocumentInCollectionNameByObjectId(client, collectionName, objectId, marshall)
	if err != nil {
		return err
	}
	return nil
}

// --- Create Utils ---

func CreateDocumentInCollection(client mongo.Client, collectionName collections.Collection, document interface{}) error {
	collection := client.Database("packagebird").Collection(collectionName.String())
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		return err
	}
	return nil
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
	var results []structures.Package
	_, err := GetObjectsFromCollectionName(client, collections.Packages.String(), &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// --- Package Set ---

func CreatePackage(client mongo.Client, pkg structures.Package) error {
	return ErrorCheck(CreateDocumentInCollection(client, collections.Packages, pkg))
}

func SetPackageByObjectId(client mongo.Client, objectId primitive.ObjectID) error {
	err := SetObjectInCollectionNameByObjectId(client, collections.Packages.String(), objectId, structures.Package{})
	if err != nil {
		return err
	}
	return nil
}

// --- Package Create

// --- Package Metadata Get ---

func GetPackageMetadataByNameAndVersion(client mongo.Client, name string, version int64) (*structures.PackageMetadata, error) {
	var filter = &bson.M{
		"$and": []bson.M{
			bson.M{
				"name": name,
			},
			bson.M{
				"version": version,
			},
		},
	}
	obj, err := GetObjectFromCollectionNameAndFilter(client, collections.PackagesMetadata.String(), structures.PackageMetadata{}, filter)
	if err != nil {
		return nil, err
	}
	return obj.(*structures.PackageMetadata), nil
}

func GetPackagesMetadata(client mongo.Client) ([]structures.PackageMetadata, error) {
	documents, err := GetObjectsFromCollectionName(client, collections.PackagesMetadata.String(), structures.PackageMetadata{})
	if err != nil {
		return nil, err
	}
	return documents.([]structures.PackageMetadata), nil
}

// --- Package Metadata Set ---

func SetPackageMetadataByObjectId(client mongo.Client, objectId primitive.ObjectID) error {
	err := SetObjectInCollectionNameByObjectId(client, collections.PackagesMetadata.String(), objectId, structures.PackageMetadata{})
	if err != nil {
		return err
	}
	return nil
}

// --- Package Metadata Create ---

func CreatePackageMetadata(client mongo.Client, packageMetadata structures.PackageMetadata) error {
	return ErrorCheck(CreateDocumentInCollection(client, collections.PackagesMetadata, packageMetadata))
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

func GetUserByName(client mongo.Client, name string) (*structures.User, error) {
	var filter = &bson.M{
		"name": name,
	}
	obj, err := GetObjectFromCollectionNameAndFilter(client, collections.Users.String(), structures.User{}, filter)
	if err != nil {
		return nil, err
	}
	return obj.(*structures.User), nil
}

func GetUsers(client mongo.Client) ([]structures.User, error) {
	documents, err := GetObjectsFromCollectionName(client, collections.Users.String(), structures.User{})
	if err != nil {
		return nil, err
	}
	return documents.([]structures.User), nil
}

// --- User Set ---

func SetUserByObjectId(client mongo.Client, objectId primitive.ObjectID) error {
	err := SetObjectInCollectionNameByObjectId(client, collections.Users.String(), objectId, structures.User{})
	if err != nil {
		return err
	}
	return nil
}

// --- User Create ---

func CreateUser(client mongo.Client, user structures.User) error {
	return ErrorCheck(CreateDocumentInCollection(client, collections.Users, user))
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

// --- Authentication Set ---

func SetAuthenticationByObjectId(client mongo.Client, objectId primitive.ObjectID) error {
	err := SetObjectInCollectionNameByObjectId(client, collections.Authentications.String(), objectId, structures.Authentication{})
	if err != nil {
		return err
	}
	return nil
}

// --- Authentication Create ---

func CreateAuthentication(client mongo.Client, authentication structures.Authentication) error {
	return ErrorCheck(CreateDocumentInCollection(client, collections.Authentications, authentication))
}

// --- Source Get ---

func GetSourceByObjectId(client mongo.Client, objectId primitive.ObjectID) (*structures.Source, error) {
	object, err := GetObjectFromCollectionNameByObjectId(client, collections.Sources.String(), objectId, &structures.Source{})
	if err != nil {
		return nil, err
	}
	return object.(*structures.Source), nil
}

func GetSources(client mongo.Client) ([]*structures.Source, error) {
	objects, err := GetObjectsFromCollectionName(client, collections.Sources.String(), &structures.Source{})
	if err != nil {
		return nil, err
	}
	return objects.([]*structures.Source), nil
}

// --- Source Set ---

func SetSourceByObjectId(client mongo.Client, objectId primitive.ObjectID) error {
	err := SetObjectInCollectionNameByObjectId(client, collections.Authentications.String(), objectId, structures.Source{})
	if err != nil {
		return err
	}
	return nil
}

// --- Source Create ---

func CreateSource(client mongo.Client, source structures.Source) error {
	err := CreateDocumentInCollection(client, collections.Sources, source)
	if err != nil {
		return err
	}
	return nil
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
	var filter = &bson.M{
		"name": name,
	}
	obj, err := GetObjectFromCollectionNameAndFilter(client, collections.Projects.String(), &structures.Project{}, filter)
	if err != nil {
		return nil, err
	}
	return obj.(*structures.Project), nil
}

func GetProjects(client mongo.Client) ([]structures.Project, error) {
	var results []structures.Project
	_, err := GetObjectsFromCollectionName(client, collections.Projects.String(), &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// --- Project Set ---

func SetProjectByObjectId(client mongo.Client, objectId primitive.ObjectID) error {
	err := SetObjectInCollectionNameByObjectId(client, collections.Projects.String(), objectId, structures.Project{})
	if err != nil {
		return err
	}
	return nil
}

// --- Project Create ---

func CreateProject(client mongo.Client, project structures.Project) error {
	return ErrorCheck(CreateDocumentInCollection(client, collections.Projects, project))
}

// --- Script Get ---

func GetScriptByObjectId(client mongo.Client, objectId primitive.ObjectID) (*structures.Script, error) {
	object, err := GetObjectFromCollectionNameByObjectId(client, collections.Scripts.String(), objectId, structures.Script{})
	if err != nil {
		return nil, err
	}
	return object.(*structures.Script), nil
}

func GetScriptByName(client mongo.Client, name string) (*structures.Script, error) {
	obj, err := GetObjectFromCollectionNameAndFilter(client, collections.Scripts.String(), filterByName(name), structures.Script{})
	if err != nil {
		return nil, err
	}
	return obj.(*structures.Script), nil
}

// --- Script Set ---

func SetScriptByObjectId(client mongo.Client, objectId primitive.ObjectID) error {
	err := SetObjectInCollectionNameByObjectId(client, collections.Scripts.String(), objectId, structures.Script{})
	if err != nil {
		return err
	}
	return nil
}

// --- Script Create ---

func CreateScript(client mongo.Client, script structures.Script) error {
	return ErrorCheck(CreateDocumentInCollection(client, collections.Scripts, script))
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
	obj, err := GetObjectFromCollectionNameAndFilter(client, collections.Graphs.String(), &structures.Graph{}, filterByName(name))
	if err != nil {
		return nil, err
	}
	return obj.(*structures.Graph), nil
}

func GetGraphs(client mongo.Client) ([]structures.Graph, error) {
	objects, err := GetObjectsFromCollectionName(client, collections.Graphs.String(), structures.Graph{})
	if err != nil {
		return nil, err
	}
	return objects.([]structures.Graph), nil
}

// --- Graph Set ---

func SetGraphByObjectId(client mongo.Client, objectId primitive.ObjectID) error {
	err := SetObjectInCollectionNameByObjectId(client, collections.Graphs.String(), objectId, structures.Graph{})
	if err != nil {
		return err
	}
	return nil
}

// --- Graph Create ---

func CreateGraph(client mongo.Client, graph structures.Graph) error {
	return ErrorCheck(CreateDocumentInCollection(client, collections.Graphs, graph))
}
