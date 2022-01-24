package NetworkInterface

import (
	"context"
	"log"
	databaseInterface "packagebird-server/src/DatabaseInterface"
	projectOperations "packagebird-server/src/NetworkInterface/ProjectOperations"
)

func (server *GRPCServer) GetProject(context context.Context, request *projectOperations.ProjectRequest) (*projectOperations.ProjectResponse, error) {

	projectName := request.GetName()
	projectDescription := request.GetDescription()

	createdProject, err := databaseInterface.NewProject(*mongoDBClientGlobal, projectName, projectDescription)
	if err != nil {
		log.Printf("Error encountered attempting to create new project...")
		return nil, err
	}

	if createdProject {
		return &projectOperations.ProjectResponse{Exist: false}, nil
	} else {
		return &projectOperations.ProjectResponse{Exist: true}, nil
	}
}
