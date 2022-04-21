package NetworkInterface

import (
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
	fileTransfer "packagebird-server/src/NetworkInterface/FileTransfer"
	buildtest "packagebird-server/src/NetworkInterface/buildtest"
	interfaces "packagebird-server/src/NetworkInterface/interfaces"
	listcontent "packagebird-server/src/NetworkInterface/listcontent"
	member_operations "packagebird-server/src/NetworkInterface/member_operations"
)

// All to-be-implemented gRPC methods must be added to this structure

// Global mongoDBClient reference
var mongoDBClientGlobal *mongo.Client

func PackagebirdServerStart(address string, mongodbClient *mongo.Client) error {

	// Creates listener at specified address with port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	} else {
		log.Printf("Listening on %v", address)
	}

	// Global reference to passed mongodb client
	mongoDBClientGlobal = mongodbClient

	// Creates new gRPC server at specified port
	server := grpc.NewServer()

	// Register passed functions with implementations, must add each set of operations
	fileTransfer.RegisterFileServiceServer(server, &interfaces.GRPCServer{})
	listcontent.RegisterListContentServicesServer(server, &interfaces.GRPCServer{})
	buildtest.RegisterBuildTestServicesServer(server, &interfaces.GRPCServer{})
	member_operations.RegisterMemberCRUDServicesServer(server, &interfaces.GRPCServer{})
	member_operations.RegisterMemberAuthenticationServer(server, &interfaces.GRPCServer{})

	log.Print("Registered gRPC methods on server...")

	// Serves gRPC operations on specified listener
	err = server.Serve(listener)
	if err != nil {
		log.Print("Error encountered in server, returning to calling method...")
		return err
	}

	// Successful operation, close server without error.
	return nil
}
