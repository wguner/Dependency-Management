/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"io/fs"
	"log"
	"os"
	"packagebird/network/grpc/services"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Creates a new packagebird project, with code contributions committed therein",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating new project directory")
		createProject(args)
	},
}

func init() {
	createCmd.AddCommand(projectCmd)
}

func createProject(projectName []string) {
	if err := os.Mkdir(projectName[0], fs.ModePerm); err != nil {
		log.Fatal("Error creating project directory")
	}
	connection, err := grpc.Dial("127.0.0.1:55051")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	client := services.NewPackagebirdServicesClient(connection)
	response, err := client.CreateProject(context.Background(), &services.PackageRequest{Name: projectName[0], Version: 0})
	if err != nil {
		panic(err)
	}
	log.Print(response.Header)
}
