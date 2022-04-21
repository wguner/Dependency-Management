package services

import (
	"context"
	"fmt"
	"github.com/arduino/go-paths-helper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"os"
	"packagebird-server/filesystem"
	"packagebird-server/global"
	"packagebird-server/mongo/exchange/accessors"
	"packagebird-server/mongo/structures"
	"path/filepath"
	"time"
)

func (server *Services) CreateProject(context context.Context, projectRequest *ProjectRequest) (*OperationResponse, error) {
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

	// Err Response
	var responseError = &OperationResponse{
		Success: false,
		Header:  "FAILED TO CREATE PROJECT",
		Message: "Failed to create project due to unknown reason",
	}

	// Success Response
	var responseSuccess = &OperationResponse{
		Success: true,
		Header:  "CREATED PROJECT",
		Message: fmt.Sprintf("Created project '%v'", projectRequest.Name),
	}

	// authentication, err := accessors.GetAuthenticationByUserObjectId()
	// authentication later
	wd, err := paths.Getwd()

	packageId := primitive.NewObjectID()
	var source = &structures.Source{
		Path:           fmt.Sprintf(filepath.ToSlash(fmt.Sprintf("%v/projects/%v/src", wd.String(), projectRequest.GetName()))),
		LastAccessedBy: time.Now(),
		ProjectName:    projectRequest.GetName(),
		ObjectId:       primitive.NewObjectID(),
	}
	var graph = &structures.Graph{
		Name:     projectRequest.GetName(),
		Version:  0,
		Package:  packageId,
		Children: []primitive.ObjectID{},
		ObjectId: primitive.NewObjectID(),
	}
	var newProject = &structures.Project{
		ObjectId:       packageId,
		Name:           projectRequest.GetName(),
		Source:         source.ObjectId,
		ProjectVersion: 0,
		PackageVersion: 0,
		Graph:          graph.ObjectId,
		Dependencies:   []primitive.ObjectID{},
	}

	err = accessors.CreateSource(*global.GlobalMongoClient, *source)
	if err != nil {
		return responseError, err
	}

	err = accessors.CreateGraph(*global.GlobalMongoClient, *graph)
	if err != nil {
		return responseError, err
	}

	err = accessors.CreateProject(*global.GlobalMongoClient, *newProject)
	if err != nil {
		return responseError, err
	}

	err = filesystem.CreateProjectSourceDirectory(projectRequest.GetName())
	if err != nil {
		return responseError, err
	}

	err = filesystem.CreateProjectSourceFromBytes(projectRequest.GetName(), []byte(""))
	if err != nil {
		return responseError, err
	}

	return responseSuccess, nil
}

func (server *Services) UploadFile(data PackagebirdServices_UploadFileServer) error {
	// Time request received
	start := time.Now()

	// Get file path
	chunk, err := data.Recv()
	if err != nil {
		return err
	}
	path := chunk.GetPath()

	// Query for source path
	project, err := accessors.GetProjectByName(*global.GlobalMongoClient, path)
	if err != nil {
		return err
	}
	source, err := accessors.GetSourceByObjectId(*global.GlobalMongoClient, project.Source)
	if err != nil {
		return err
	}
	path = source.Path

	// Check if file exist
	_, err = os.Stat(path)
	if err != nil {
		return err
	}

	// Write file contents
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	for {
		chunk, err := data.Recv()
		if (chunk == nil) || (len(chunk.GetChunk()) == 0) {
			break
		}
		if err != nil {
			return err
		}

		_, err = file.Write(chunk.GetChunk())
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	// Notify client
	data.SendAndClose(&OperationResponse{
		Header:  "SUCCESSFULLY UPLOADED FILE",
		Success: true,
		Message: fmt.Sprintf("Uploaded file to path '%v' successfully in %v minutes", path, time.Now().Sub(start).Minutes()),
	})
	return nil
}

func (server *Services) DownloadFile(request *DownloadRequest, data PackagebirdServices_DownloadFileServer) error {
	path := request.GetPath()

	// Get project by name
	project, err := accessors.GetProjectByName(*global.GlobalMongoClient, path)
	if err != nil {
		return err
	}

	// Get project source
	source, err := accessors.GetSourceByObjectId(*global.GlobalMongoClient, project.Source)
	if err != nil {
		return err
	}

	// Get project file
	path = source.Path

	// Open file
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// Write file to download
	buffer := make([]byte, 64*1024)
	for {
		bytes, err := file.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		chunk := &File{
			Content: &File_Chunk{
				Chunk: buffer[:bytes],
			},
		}

		err = data.Send(chunk)
		if err != nil {
			return err
		}
	}

	return nil
}