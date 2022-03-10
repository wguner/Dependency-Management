package structures

import "go.mongodb.org/mongo-driver/mongo"

var GlobalMongoClient *mongo.Client

const (
	PACKAGEPATH = "NOTHING"
)

type Package struct {
	Name           string          `bson:"name"`
	Description    string          `bson:"description"`
	UUID           string          `bson:"uuid"`
	Authors        []string        `bson:"authors"`
	Version        int64           `bson:"version"`
	SourceFile     string          `bson:"source_file"`  // Path to source code, binaries on disk for server
	Dependencies   []string        `bson:"dependencies"` // Recursive reference to other packages
	Graph          DependencyGraph `bson:"graph,inline"`
	CreatedOn      string          `bson:"created_on"`
	LastAccessedOn string          `bson:"last_accessed_on"`
	BuildFile      string          `bson:"build_file"`
	TestFile       string          `bson:"test_file"`
}

type Project struct {
	Name          string   `bson:"name"`
	Description   string   `bson:"description"`
	UUID          string   `bson:"uuid"`
	LatestVersion int64    `bson:"latest_version"`
	Packages      []string `bson:"packages"`
	SourceFile    string   `bson:"source_file"`
	BuildFile     string   `bson:"build_file"`
	TestFile      string   `bson:"test_file"`
	CreatedOn     string   `bson:"created_on"`
	LastSyncedOn  string   `bson:"last_synced_on"`
	LastSyncedBy  string   `bson:"last_synced_by"`
}

type Team struct {
	TeamName string    `bson:"teamname"`
	Members  []Member  `bson:"members"`
	Projects []Project `bson:"projects"`
}

type Member struct {
	Name     string `bson:"name"`
	Password string `bson:"password"`
	Employed bool   `bson:"is_employed"`
	Level    string `bson:"level"`
}

type DependencyGraph struct {
	Dependencies []Package `bson:"dependencies"`
}
