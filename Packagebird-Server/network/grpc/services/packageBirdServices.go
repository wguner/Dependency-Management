package services

import (
	"context"
	"fmt"
	"github.com/arduino/go-paths-helper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"os"
	"packagebird-server/filesystem"
	"packagebird-server/global"
	"packagebird-server/mongo/exchange/accessors"
	"packagebird-server/mongo/structures"
	"path/filepath"
	"strings"
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
	split := strings.Split(path, "/")
	name := split[len(split)-4]

	// Query for source path
	project, err := accessors.GetProjectByName(*global.GlobalMongoClient, name)
	if err != nil {
		return err
	}
	pkgVersion := project.PackageVersion - 1
	pkg, err := accessors.GetPackageByNameAndVersion(*global.GlobalMongoClient, name, pkgVersion)
	source, err := accessors.GetSourceByObjectId(*global.GlobalMongoClient, pkg.Source)
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

func (server *Services) GetProjects(context context.Context, blank *Blank) (*ProjectList, error) {
	projects, err := accessors.GetProjects(*global.GlobalMongoClient)
	if err != nil {
		return nil, err
	}

	var names []*ProjectName
	for _, ele := range projects {
		names = append(names, &ProjectName{Name: ele.Name})
	}

	return &ProjectList{Names: names}, nil
}

func (server *Services) CreatePackage(context context.Context, request *PackageRequest) (*OperationResponse, error) {
	var response = &OperationResponse{
		Success: false,
		Header:  "FAILED TO CREATE PACKAGE",
	}

	project, err := accessors.GetProjectByName(*global.GlobalMongoClient, request.GetName())
	if err != nil && err != mongo.ErrNoDocuments {
		response.Header = "FAILED TO RETRIEVE PROJECT"
		response.Message = fmt.Sprintf("Failed to retrieve project by name '%v'", request.GetName())
		return response, err
	}

	// TODO: authentication, err := accessors.GetAuthenticationByUserObjectId()
	// TODO: authentication later

	workingDirectory, err := paths.Getwd()

	var source = &structures.Source{
		Path:           fmt.Sprintf(filepath.ToSlash(fmt.Sprintf("%v/packages/%v/version/%v/src", workingDirectory.String(), request.GetName(), project.PackageVersion))),
		LastAccessedBy: time.Now(),
		PackageName:    request.GetName(),
		ObjectId:       primitive.NewObjectID(),
	}
	var graph = &structures.Graph{
		ObjectId: primitive.NewObjectID(),
		Name:     project.Name,
		Version:  project.PackageVersion,
		Package:  primitive.NewObjectID(),
		Children: project.Dependencies,
	}
	var pkg = &structures.Package{
		Name:     request.GetName(),
		Version:  project.PackageVersion,
		Source:   source.ObjectId,
		Graph:    graph.ObjectId,
		Scripts:  []primitive.ObjectID{},
		ObjectId: graph.Package,
	}

	err = accessors.CreateSource(*global.GlobalMongoClient, *source)
	if err != nil {
		return response, err
	}

	err = accessors.CreateGraph(*global.GlobalMongoClient, *graph)
	if err != nil {
		return response, err
	}

	err = accessors.CreatePackage(*global.GlobalMongoClient, *pkg)
	if err != nil {
		return response, err
	}

	err = filesystem.CreatePackageSourceDirectory(pkg.Name, pkg.Version)
	if err != nil {
		return response, err
	}

	// Update package version attached to project
	var update = &bson.D{
		{"$set", &bson.D{
			{"packageVersion", project.PackageVersion + 1},
		}},
	}
	err = accessors.SetProjectByObjectId(*global.GlobalMongoClient, project.ObjectId, update)

	response.Success = true
	response.Header = "SUCCESSFULLY CREATED PACKAGE TEMPLATE"
	response.Message = source.Path
	response = &OperationResponse{
		Success: true,
		Header:  "SUCCESSFULLY CREATED PACKAGE TEMPLATE",
		Message: source.Path,
	}
	return response, nil
}

func (server *Services) GetPackages(ctx context.Context, blank *Blank) (*PackageList, error) {
	pkgs, err := accessors.GetPackages(*global.GlobalMongoClient)
	if err != nil {
		return nil, err
	}
	var names []*PackageName
	for _, ele := range pkgs {
		names = append(names, &PackageName{Name: ele.Name})
	}
	return &PackageList{Names: names}, nil
}
