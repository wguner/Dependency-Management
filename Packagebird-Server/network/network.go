package network

import (
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"net"
	"packagebird-server/network/grpc/services"
)

func StartServicesServer(address string, mongodb *mongo.Client) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	services.RegisterPackagebirdServicesServer(server, &services.Services{})

	err = server.Serve(listen)
	if err != nil {
		return err
	}

	return nil
}
