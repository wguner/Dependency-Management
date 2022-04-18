package structures

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

// Graph node within dependency graph, references package
type Graph struct {
	Name     string          `bson:"name"`
	Version  int64           `bson:"version"`
	Package  bson.ObjectId   `bson:"package"`
	Children []bson.ObjectId `bson:"children"`
}

// Package representation of a module of source code
type Package struct {
	Name    string          `bson:"name"`
	Version int64           `bson:"version"`
	Source  bson.ObjectId   `bson:"sourceId"`
	Graph   bson.ObjectId   `bson:"graphId"`
	Scripts []bson.ObjectId `bson:"scripts"`
}

// PackageMetadata metadata attached to a package
type PackageMetadata struct {
	Package          bson.ObjectId `bson:"packageId"`
	NumberDownloads  int64         `bson:"numberDownloads"`
	LastDownloaded   time.Time     `bson:"lastDownloaded"`
	LastDownloadedBy bson.ObjectId `bson:"LastDownloadedBy"`
}

// User in system
type User struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

// Authentication entry for particular user in system
type Authentication struct {
	User    bson.ObjectId   `bson:"userId"`
	Project []bson.ObjectId `bson:"projectIds"`
	IsAdmin bool            `bson:"isAdmin"`
}

// Source representation of the actual backing file for a project or a package
type Source struct {
	Path           string    `bson:"path"`
	LastAccessedBy time.Time `bson:"LastAccessedBy"`
}

// Project a package under active mutable development
type Project struct {
	Name           string          `bson:"name"`
	Source         bson.ObjectId   `bson:"sourceId"`
	ProjectVersion int64           `bson:"projectVersion"`
	PackageVersion int64           `bson:"packageVersion"`
	Graph          bson.ObjectId   `bson:"graphId"`
	Dependencies   []bson.ObjectId `bson:"packages"`
}

// Script a runnable script attached to a package
type Script struct {
	Name        string          `bson:"name"`
	Description string          `bson:"description"`
	Body        string          `bson:"body"`
	Package     bson.ObjectId   `bson:"package"`
	Packages    []bson.ObjectId `bson:"packages"`
}
