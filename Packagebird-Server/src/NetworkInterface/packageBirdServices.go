package NetworkInterface

import (
	"context"
	"fmt"
	"github.com/arduino/go-paths-helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"packagebird-server/global"
	"packagebird-server/mongo/exchange/accessors"
	"packagebird-server/mongo/structures"
	"path/filepath"
	"time"
)

func (server *GRPCServer) CreateProject(context context.Context, projectRequest *ProjectRequest) (*OperationResponse, error) {
	proj, err := accessors.GetProjectByName(*global.GlobalMongoClient, projectRequest.GetName())
	if err != nil && err != mongo.ErrNoDocuments {
		var response = &OperationResponse{
			Success: false,
			Header:  "FAILED TO RETRIEVE PROJECT",
			Message: "Unknown failure attempting to find project if already exist.",
		}
		return response, err
	}
	if proj != nil {
		var response = &OperationResponse{
			Success: false,
			Header:  "PROJECT ALREADY EXIST",
			Message: "Project by that name already exist. Please pick another project name.",
		}
		return response, nil
	}

	// authentication, err := accessors.GetAuthenticationByUserObjectId()
	// authentication later
	wd, err := paths.Getwd()

	var source = &structures.Source{
		Path:           fmt.Sprintf(filepath.ToSlash(fmt.Sprintf("%v/projects/%v/%v", wd.String(), projectRequest.GetName(), 0))),
		LastAccessedBy: time.Now(),
		ProjectName:    projectRequest.GetName(),
	}

	err = accessors.CreateSource(*global.GlobalMongoClient, *source)
	if err != nil {
		return nil, err
	}

	arr, err := accessors.GetSources(*global.GlobalMongoClient)
	f := func(e structures.Source) bool {
		return e.ProjectName == projectRequest.GetName()
	}
	arr = accessors.Filter(arr, f)
	source = &arr[0]

	var graph = &structures.Graph{
		Name:     projectRequest.GetName(),
		Version:  0,
		Package:  primitive.ObjectID{},
		Children: []primitive.ObjectID{},
	}

	err = accessors.CreateGraph(*global.GlobalMongoClient, *graph)

	graph, err = accessors.GetGraphByName(*global.GlobalMongoClient, projectRequest.GetName())
	if err != nil {
		return nil, err
	}

	var newProject = &structures.Project{
		Name:           projectRequest.GetName(),
		Source:         source.ObjectId,
		ProjectVersion: 0,
		PackageVersion: 0,
		Graph:          graph.ObjectId,
		Dependencies:   []primitive.ObjectID{},
	}

	err = accessors.CreateProject(*global.GlobalMongoClient, *newProject)
	if err != nil {
		var response = &OperationResponse{
			Success: false,
			Header:  "FAILED TO CREATE PROJECT",
			Message: "Failed to create project due to unknown reason",
		}
		return response, err
	}
	var response = &OperationResponse{
		Success: true,
		Header:  "CREATED PROJECT",
		Message: fmt.Sprintf("Created project '%v'", newProject.Name),
	}
	return response, nil
}
