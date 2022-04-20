package NetworkInterface

import (
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
	fileTransfer "packagebird-server/src/NetworkInterface/FileTransfer"
	serverUtils "packagebird-server/src/NetworkInterface/ServerUtils"
	buildtest "packagebird-server/src/NetworkInterface/buildtest"
	listcontent "packagebird-server/src/NetworkInterface/listcontent"
	member_operations "packagebird-server/src/NetworkInterface/member_operations"
)

// All to-be-implemented gRPC methods must be added to this structure
type GRPCServer struct {
	fileTransfer.UnimplementedFileServiceServer
	serverUtils.UnimplementedServerUtilsServicesServer
	listcontent.UnimplementedListContentServicesServer
	buildtest.UnimplementedBuildTestServicesServer
	member_operations.UnimplementedMemberCRUDServicesServer
	member_operations.UnimplementedMemberAuthenticationServer
	UnimplementedPackagebirdServicesServer
}

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
	RegisterPackagebirdServicesServer(server, &GRPCServer{})
	fileTransfer.RegisterFileServiceServer(server, &GRPCServer{})
	serverUtils.RegisterServerUtilsServicesServer(server, &GRPCServer{})
	listcontent.RegisterListContentServicesServer(server, &GRPCServer{})
	buildtest.RegisterBuildTestServicesServer(server, &GRPCServer{})
	member_operations.RegisterMemberCRUDServicesServer(server, &GRPCServer{})
	member_operations.RegisterMemberAuthenticationServer(server, &GRPCServer{})

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
