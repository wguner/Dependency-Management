/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"packagebird/network/grpc/services"

	"github.com/spf13/cobra"
)

// projectsCmd represents the projects command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(getProjects())
	},
}

func init() {
	getCmd.AddCommand(projectsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// projectsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getProjects() error {
	connection, err := grpc.Dial(GetServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil
	}
	defer connection.Close()

	client := services.NewPackagebirdServicesClient(connection)
	response, err := client.GetProjects(context.Background(), &services.Blank{})
	if err != nil {
		return err
	}

	for _, ele := range response.GetNames() {
		fmt.Printf("%v\n", ele.GetName())
	}
	return nil
}
