package NetworkInterface

import (
	"log"
	"net"

	fileTransfer "packagebird-server/src/NetworkInterface/FileTransfer"
	packageOperations "packagebird-server/src/NetworkInterface/PackageOperations"
	projectOperations "packagebird-server/src/NetworkInterface/ProjectOperations"
	serverUtils "packagebird-server/src/NetworkInterface/ServerUtils"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

// All to-be-implemented gRPC methods must be added to this structure
type GRPCServer struct {
	packageOperations.UnimplementedPackageOperationServicesServer
	projectOperations.UnimplementedProjectOperationServicesServer
	fileTransfer.UnimplementedFileServiceServer
	serverUtils.UnimplementedServerUtilsServicesServer
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
	packageOperations.RegisterPackageOperationServicesServer(server, &GRPCServer{})
	projectOperations.RegisterProjectOperationServicesServer(server, &GRPCServer{})
	fileTransfer.RegisterFileServiceServer(server, &GRPCServer{})
	serverUtils.RegisterServerUtilsServicesServer(server, &GRPCServer{})

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
