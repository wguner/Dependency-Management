package NetworkInterface

import (
	"context"
	"fmt"
	database_interface "packagebird-server/src/DatabaseInterface"
	member_operations "packagebird-server/src/NetworkInterface/member_operations"
	structures "packagebird-server/src/structures"

	"go.mongodb.org/mongo-driver/mongo"
)

// Add a member to the server
func (server *GRPCServer) AddMember(context context.Context, request *member_operations.MemberRequest) (*member_operations.MemberResponse, error) {
	// Map request to local structure
	member := structures.Member{
		Name:     request.Name,
		Password: request.Password,
		Level:    request.Level,
		Employed: true,
	}

	// Create new member in database
	if err := database_interface.CreateMember(*mongoDBClientGlobal, structures.Member(member)); err != nil {
		// If member already exist or error creating member, cancel and notify client
		return &member_operations.MemberResponse{Response: fmt.Sprintf("Could not create new member with name '%v' in database. Member may already exist in database.", member.Name), WasSuccessful: false}, err
	} else {
		// If member registered successfully, notify client
		return &member_operations.MemberResponse{Response: fmt.Sprintf("Successfully registered new member with '%v' in database.", member.Name), WasSuccessful: true}, nil
	}
}

// Remove a member from the server
func (server *GRPCServer) RemoveMember(context context.Context, request *member_operations.MemberRequest) (*member_operations.MemberResponse, error) {
	if err := database_interface.RemoveMember(*mongoDBClientGlobal, request.Name); err != nil {
		return &member_operations.MemberResponse{Response: fmt.Sprintf("Could not remove member with name '%v' in database. Member may not exist in database.", request.Name), WasSuccessful: false}, err
	}
	return &member_operations.MemberResponse{Response: fmt.Sprintf("Removed member with name '%v' in database.", request.Name), WasSuccessful: true}, nil
}

// Authenticate a client's request
func (server *GRPCServer) AuthenticateMember(context context.Context, request *member_operations.MemberAuthenticationRequest) (*member_operations.MemberAuthenticationResponse, error) {
	// If this is set, check if client is admin user
	if request.GetIsLevel() {
		// Get admin user, return problem response if pressent
		admin, err := database_interface.GetAdminMember(*mongoDBClientGlobal)
		if err != nil {
			response := &member_operations.MemberAuthenticationResponse{Response: "Error attempting to retrieve admin user for server", IsAuthenticated: false}
			return response, err
		}

		// If name and password match with admin, authenticate as admin and return
		if admin.Name == request.Name && admin.Password == request.Password {
			response := &member_operations.MemberAuthenticationResponse{Response: "Successfully authenticated administrative user", IsAuthenticated: true}
			return response, nil
		}
	}

	// Else authenticate as normal user
	member, err := database_interface.GetMember(*mongoDBClientGlobal, request.Name)
	// If no member encountered, relay information back to user
	if err == mongo.ErrNoDocuments {
		response := &member_operations.MemberAuthenticationResponse{Response: fmt.Sprintf("Error; no member with name '%v' registered on server.", request.Name), IsAuthenticated: false}
		return response, nil
	}

	// Else, check if passwords match
	if member.Password == request.Password {
		// If successful, notify client
		return &member_operations.MemberAuthenticationResponse{Response: fmt.Sprintf("Successfully authenticated member with name '%v'", request.Name), IsAuthenticated: true}, nil
	} else {
		// If failure, notify client
		return &member_operations.MemberAuthenticationResponse{Response: fmt.Sprintf("Failed to authenticate member with name '%v'", request.Name), IsAuthenticated: false}, nil
	}
}
