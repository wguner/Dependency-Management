package NetworkInterface

import (
	"log"
	"net"

	fileTransfer "packagebird-server/src/NetworkInterface/FileTransfer"
	packageOperations "packagebird-server/src/NetworkInterface/PackageOperations"

	"google.golang.org/grpc"
)

// All to-be-implemented gRPC methods must be added to this structure
type grpcServer struct {
	packageOperations.UnimplementedPackageOperationServicesServer
	fileTransfer.UnimplementedFileServiceServer
}

func PackagebirdServerStart(address string) error {

	// Creates listener at specified address with port
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	} else {
		log.Printf("Listening on %v", address)
	}

	// Creates new gRPC server at specified port
	server := grpc.NewServer()

	// Register passed functions with implementations, must add each set of operations
	packageOperations.RegisterPackageOperationServicesServer(server, &grpcServer{})
	fileTransfer.RegisterFileServiceServer(server, &grpcServer{})

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
