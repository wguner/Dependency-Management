package main

import (
	protobuf_grpc "Server/FileTransfer"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

const (
	ADDRESS     = "127.0.0.1"
	PORT        = 50051
	BUFFER_SIZE = 64 * 1024
)

type Server struct {
	protobuf_grpc.UnimplementedFileServiceServer
}

func (s *Server) Download(request *protobuf_grpc.Request, fileStream protobuf_grpc.FileService_DownloadServer) error {

	file, err := os.Open(request.GetBody())
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

	buffer := make([]byte, BUFFER_SIZE)

	for {
		bytesRead, err := file.Read(buffer)

		// Oddly written way of intercepting EOF
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}

		file := &protobuf_grpc.File{
			Content: &protobuf_grpc.File_Chunk{
				Chunk: buffer[:bytesRead],
			},
		}

		err = fileStream.Send(file)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (s *Server) Upload(fileStream protobuf_grpc.FileService_UploadServer) error {

	file, err := os.Create("tmp/temp.bin")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	for {
		chunk, err := fileStream.Recv()
		if (chunk == nil) || (len(chunk.GetChunk()) == 0) {
			break
		}

		if err != nil {
			log.Fatal(err)
			return err
		}

		_, err = file.Write(chunk.GetChunk())
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
			return err
		}
	}

	message := &protobuf_grpc.Response{
		Body: "File uploaded successfully!",
	}

	fileStream.SendAndClose(message)

	return nil
}

func (s *Server) NameFile(context context.Context, request *protobuf_grpc.Request) (*protobuf_grpc.Response, error) {
	oldFileName := "tmp/temp.bin"
	newFileName := fmt.Sprintf("tmp/%v", request.GetBody())

	err := os.Rename(oldFileName, newFileName)
	if err != nil {
		log.Fatal(err)

		response := &protobuf_grpc.Response{
			Header: "File Rename",
			Body:   "Failed to rename file!",
		}

		return response, err
	}

	response := &protobuf_grpc.Response{
		Header: "File Rename",
		Body:   "Succesfully renamed file!",
	}

	return response, nil
}

func main() {
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", ADDRESS, PORT))
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}

	server := grpc.NewServer()
	protobuf_grpc.RegisterFileServiceServer(server, &Server{})
	log.Printf("Server listening at address %v", listener.Addr())

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error serving: %v", err)
	}
}
