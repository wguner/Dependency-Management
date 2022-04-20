/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/fs"
	"log"
	"os"
	services "packagebird/network/grpc/services"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Creates a new packagebird project, with code contributions committed therein",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating new project directory")
		createProject(args[0])
	},
}

func init() {
	createCmd.AddCommand(projectCmd)
}

func createProject(projectName string) {
	if err := os.Mkdir(projectName, fs.ModePerm); err != nil {
		log.Fatal("Error creating project directory")
	}
	connection, err := grpc.Dial("127.0.0.1:55051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	client := services.NewPackagebirdServicesClient(connection)
	response, err := client.CreateProject(context.Background(), &services.ProjectRequest{Name: projectName})
	if err != nil {
		panic(err)
	}
	log.Print(response.Header)
}
