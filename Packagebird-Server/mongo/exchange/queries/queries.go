package queries

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"packagebird-server/mongo/exchange/accessors"
	"packagebird-server/mongo/structures"
)

// GetAllPackageDependencies returns all dependencies of a package
func GetAllPackageDependenciesByNameAndVersion(client mongo.Client, name string, version int64) ([]*structures.Package, error) {

	// Retrieve package by name and version
	dep, err := accessors.GetPackageByNameAndVersion(client, name, version)
	if err != nil {
		return nil, err
	}

	// Retrieve graph attached to package
	graph, err := accessors.GetGraphByObjectId(client, dep.Graph)
	if err != nil {
		return nil, err
	}

	// Recursively traverse and collect nodes from graph
	var nodes []*structures.Graph
	err = GetGraph(client, *graph, nodes)
	if err != nil {
		return nil, err
	}

	// Potential Duplicates present in array
	children := RemoveDuplicates(nodes)

	// Collect all packages present in array
	var packages []*structures.Package
	for _, child := range children {
		packageId := child.Package
		pck, err := accessors.GetPackageByObjectId(client, packageId)
		if err != nil {
			return nil, err
		}
		packages = append(packages, pck)
	}

	return packages, nil
}

func GetAllPackageDependenciesByObjectId(client mongo.Client, objectId primitive.ObjectID) ([]*structures.Package, error) {
	pkg, err := accessors.GetPackageByObjectId(client, objectId)
	if err != nil {
		return nil, err
	}

	graph, err := accessors.GetGraphByObjectId(client, pkg.Graph)
	if err != nil {
		return nil, err
	}

	var deps []*structures.Graph
	err = GetGraph(client, *graph, deps)
	if err != nil {
		return nil, err
	}

	children := RemoveDuplicates(deps)

	var pkgs []*structures.Package
	for _, child := range children {
		pkgId := child.Package
		pkg, err := accessors.GetPackageByObjectId(client, pkgId)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, pkg)
	}
	return pkgs, nil
}

// GetGraph recursively traverses a graph and returns the associated nodes
func GetGraph(client mongo.Client, graph structures.Graph, nodes []*structures.Graph) error {
	nodes = append(nodes, &graph)
	for _, child := range graph.Children {
		graph, err := accessors.GetGraphByObjectId(client, child)
		if err != nil {
			return err
		}

		err = GetGraph(client, *graph, nodes)
		if err != nil {
			return err
		}
	}
	return nil
}

// RemoveDuplicates removes all duplicate items from an array
func RemoveDuplicates[T any](array []T) []T {
	keys := make(map[*T]bool)
	var list []T
	for _, item := range array {
		if _, value := keys[&item]; !value {
			keys[&item] = true
			list = append(list, item)
		}
	}
	return list
}

// GetAuthenticatedProjectsByUserName retrieves all projects attached to authentication by user name
func GetAuthenticatedProjectsByUserName(client mongo.Client, name string) ([]structures.Project, error) {
	// Get user
	user, err := accessors.GetUserByName(client, name)
	if err != nil {
		return nil, err
	}

	// Get authentication
	auth, err := accessors.GetAuthenticationByUserObjectId(client, user.Authentication)

	// Get projects attached to authentication
	var projects []structures.Project
	for _, project := range auth.Project {
		proj, err := accessors.GetProjectByObjectId(client, project)
		if err != nil {
			return nil, err
		}
		projects = append(projects, *proj)
	}
	return projects, nil
}

func GetPackageScriptsByNameAndVersion(client mongo.Client, name string, version int64) ([]structures.Script, error) {
	// Get Project
	pkg, err := accessors.GetPackageByNameAndVersion(client, name, version)
	if err != nil {
		return nil, err
	}

	// Get Scripts
	var scripts []structures.Script
	for _, scriptId := range pkg.Scripts {
		script, err := accessors.GetScriptByObjectId(client, scriptId)
		if err != nil {
			return nil, err
		}
		scripts = append(scripts, *script)
	}
	return scripts, nil
}

func GetAllDependencySources(client mongo.Client, dependencies []*structures.Package) ([]*structures.Source, error) {
	var sources []*structures.Source
	for _, dep := range dependencies {
		source, err := accessors.GetSourceByObjectId(client, dep.Source)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}
	return sources, nil
}

func GetAllSourcePaths(sources []*structures.Source) ([]string, error) {
	var paths []string
	for _, source := range sources {
		paths = append(paths, source.Path)
	}
	return paths, nil
}

func GetAllPackageDependencySourcePathsByNameAndVersion(client mongo.Client, name string, version int64) ([]string, error) {
	dependencies, err := GetAllPackageDependenciesByNameAndVersion(client, name, version)
	if err != nil {
		return nil, err
	}

	source, err := GetAllDependencySources(client, dependencies)
	if err != nil {
		return nil, err
	}

	paths, err := GetAllSourcePaths(source)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

func GetAllPackageDependencySourcePathsByObjectId(client mongo.Client, objectId primitive.ObjectID) ([]string, error) {
	deps, err := GetAllPackageDependenciesByObjectId(client, objectId)
	if err != nil {
		return nil, err
	}

	src, err := GetAllDependencySources(client, deps)
	if err != nil {
		return nil, err
	}

	paths, err := GetAllSourcePaths(src)
	if err != nil {
		return nil, err
	}
	return paths, nil
}
