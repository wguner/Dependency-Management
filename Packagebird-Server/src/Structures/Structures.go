package structures

type Package struct {
	Name         string          `bson:"name"`
	Description  string          `bson:"description"`
	UUID         string          `bson:"uuid"`
	Authors      []string        `bson:"authors,inline"`
	Version      int64           `bson:"version"`
	Source       string          `bson:"source"`              // Path to source code, binaries on disk for server
	Dependencies []Package       `bson:"dependencies,inline"` // Recursive reference to other packages
	Graph        DependencyGraph `bson:"graph,inline"`
}

type Project struct {
	Name          string    `bson:"name"`
	Description   string    `bson:"description"`
	UUID          string    `bson:"uuid"`
	LatestVersion int64     `bson:"latest_version"`
	Packages      []Package `bson:"packages,inline"`
}

type Team struct {
	TeamName string   `bson:"teamname"`
	Members  []Member `bson:"members,inline"`
}

type Member struct {
	Name     string `bson:"name"`
	Password string `bson:"password"`
	Employed bool   `bson:"is_employed"`
}

type DependencyGraph struct {
	Dependencies []Package `bson:"dependencies,inline"`
}
