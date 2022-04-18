package mongo

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
	"packagebird-server/mongo/structures"
)

// --- Utils ---

func GetDocumentFromCollectionNameByObjectId(client mongo.Client, collectionName string, objectId bson.ObjectId) (*mongo.SingleResult, error) {
	// TODO: To be filled
	return nil, nil
}

func GetDocumentsFromCollectionName(client mongo.Client, collectionName string) (*mongo.Cursor, error) {
	// TODO: To be filled
	return nil, nil
}

// --- Package Get ---

func GetPackageByObjectId(client mongo.Client, objectId bson.ObjectId) (*structures.Project, error) {
	// TODO: To be filled
	return nil, nil
}

func GetPackageByNameAndVersion(client mongo.Client, name string, version int64) (*structures.Package, error) {
	// TODO: To be filled
	return nil, nil
}

func GetPackagesByName(client mongo.Client, name string) ([]structures.Package, error) {
	// TODO: To be filled
	return nil, nil
}

func GetPackages(client mongo.Client) ([]structures.Package, error) {
	// TODO: To be filled
	return nil, nil
}

// --- Package Metadata Get ---

func GetPackageMetadataByNameAndVersion(client mongo.Client, name string, version int64) (*structures.PackageMetadata, error) {
	// TODO: To be filled
	return nil, nil
}

func GetPackagesMetadata(client mongo.Client) ([]structures.PackageMetadata, error) {
	// TODO: To be filled
	return nil, nil
}

// --- User Get ---

func GetUserByObjectId(client mongo.Client, objectId bson.ObjectId) (*structures.User, error) {
	// TODO: To be filled
	return nil, nil
}

func GetUserByName(client mongo.Client, name string) ([]structures.User, error) {
	// TODO: To be filled
	return nil, nil
}

func GetUsers(client mongo.Client) ([]structures.User, error) {
	// TODO: To be filled
	return nil, nil
}

// --- Authentication Get ---

func GetAuthenticationByUserObjectId(client mongo.Client, objectId bson.ObjectId) (*structures.Authentication, error) {
	// TODO: To be filled
	return nil, nil
}

func GetAuthentications(client mongo.Client) ([]structures.Authentication, error) {
	// TODO: To be filled
	return nil, nil
}

// --- Source Get ---

func GetSourceByObjectId(client mongo.Client, objectId bson.ObjectId) (*structures.Source, error) {
	// TODO: To be filled
	return nil, nil
}

func GetSources(client mongo.Client) ([]structures.Source, error) {
	// TODO: To be filled
	return nil, nil
}

// --- Project Get ---

func GetProjectByObjectId(client mongo.Client) (*structures.Project, error) {
	// TODO: To be filled
	return nil, nil
}

func GetProjectByName(client mongo.Client, name string) (*structures.Project, error) {
	// TODO: To be filled
	return nil, nil
}

func GetProjects(client mongo.Client) ([]structures.Project, error) {
	// TODO: To be filled
	return nil, nil
}

// --- Script Get ---

func GetScriptByObjectId(client mongo.Client) (*structures.Script, error) {
	// TODO: To be filled
	return nil, nil
}

func GetScriptByName(client mongo.Client) (*structures.Script, error) {
	// TODO: To be filled
	return nil, nil
}

// --- Graph Get ---

func GetGraphByObjectId(client mongo.Client, objectId bson.ObjectId) (*structures.Graph, error) {
	// TODO: To be filled
	return nil, nil
}

func GetGraphByName(client mongo.Client, name string) (*structures.Graph, error) {
	// TODO: To be filled
	return nil, nil
}

func GetGraphs(client mongo.Client) (*structures.Graph, error) {
	// TODO: To be filled
	return nil, nil
}
