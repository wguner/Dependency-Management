package structures

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Graph node within dependency graph, references package
type Graph struct {
	Name     string               `bson:"name"`
	Version  int64                `bson:"version"`
	Package  primitive.ObjectID   `bson:"package"`
	Children []primitive.ObjectID `bson:"children"`
}

// Package representation of a module of source code
type Package struct {
	Name    string               `bson:"name"`
	Version int64                `bson:"version"`
	Source  primitive.ObjectID   `bson:"sourceId"`
	Graph   primitive.ObjectID   `bson:"graphId"`
	Scripts []primitive.ObjectID `bson:"scripts"`
}

// PackageMetadata metadata attached to a package
type PackageMetadata struct {
	Package          primitive.ObjectID `bson:"packageId"`
	NumberDownloads  int64              `bson:"numberDownloads"`
	LastDownloaded   time.Time          `bson:"lastDownloaded"`
	LastDownloadedBy primitive.ObjectID `bson:"LastDownloadedBy"`
}

// User in system
type User struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

// Authentication entry for particular user in system
type Authentication struct {
	User    primitive.ObjectID   `bson:"userId"`
	Project []primitive.ObjectID `bson:"projectIds"`
	IsAdmin bool                 `bson:"isAdmin"`
}

// Source representation of the actual backing file for a project or a package
type Source struct {
	Path           string    `bson:"path"`
	LastAccessedBy time.Time `bson:"LastAccessedBy"`
}

// Project a package under active mutable development
type Project struct {
	Name           string               `bson:"name"`
	Source         primitive.ObjectID   `bson:"sourceId"`
	ProjectVersion int64                `bson:"projectVersion"`
	PackageVersion int64                `bson:"packageVersion"`
	Graph          primitive.ObjectID   `bson:"graphId"`
	Dependencies   []primitive.ObjectID `bson:"packages"`
}

// Script a runnable script attached to a package
type Script struct {
	Name        string               `bson:"name"`
	Description string               `bson:"description"`
	Body        string               `bson:"body"`
	Package     primitive.ObjectID   `bson:"package"`
	Packages    []primitive.ObjectID `bson:"packages"`
}
