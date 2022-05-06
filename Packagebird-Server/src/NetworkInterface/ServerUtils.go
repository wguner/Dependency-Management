package NetworkInterface

import (
	serverUtils "packagebird-server/src/NetworkInterface/ServerUtils"
	"context"
	"log"
)

// Catches ping from client to server, logs and responds

func (server *GRPCServer) Ping(context context.Context, request *serverUtils.ClientInfo) (*serverUtils.ServerInfo, error) {
	client := request.GetBody()
	log.Printf("Recieved a ping from a client: %v...", client)


	response := &serverUtils.ServerInfo{Body: "ACK"}
	return response, nil
}

