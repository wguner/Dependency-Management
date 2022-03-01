package NetworkInterface

import (
	"context"
	"log"
	databaseinterface "packagebird-server/src/DatabaseInterface"
	listcontent "packagebird-server/src/NetworkInterface/listcontent"
	structures "packagebird-server/src/structures"
)

// List Contents of the Server specified in the request
func (server *GRPCServer) GetContent(context context.Context, ContentRequest *listcontent.ContentRequest) (*listcontent.ContentResponse, error) {
	// var ContentResponse *listcontent.ContentResponse
	listProjects, listPackages, listMembers := []string{}, []string{}, []string{}
	var cErr error

	// Get list of Projects
	if ContentRequest.GetListProjects() {
		results, err := databaseinterface.GetProjects(*structures.GlobalMongoClient)
		if err != nil {
			cErr = err
		}
		if results != nil {
			listProjects = results
		}
	}

	// Get list of Packages
	if ContentRequest.GetListPackages() {
		results, err := databaseinterface.GetPackages(*structures.GlobalMongoClient)
		if err != nil {
			cErr = err
		}
		if results != nil {
			listPackages = results
		}
	}

	// Get list of Members
	if ContentRequest.GetListMembers() {
		results, err := databaseinterface.GetMembers(*structures.GlobalMongoClient)
		if err != nil {
			cErr = err
		}
		if results != nil {
			listMembers = results
		}
	}

	// Response
	if cErr != nil {
		log.Printf("Error retreiving response from database\nError:%v", cErr)
		return nil, cErr
	}
	ContentResponse := &listcontent.ContentResponse{Projects: listProjects, Packages: listPackages, Members: listMembers}
	return ContentResponse, nil
}
