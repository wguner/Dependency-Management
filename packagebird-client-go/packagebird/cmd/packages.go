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

// packagesCmd represents the packages command
var packagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		getPackages()
	},
}

func init() {
	getCmd.AddCommand(packagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getPackages() error {
	connection, err := grpc.Dial("127.0.0.1:55051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer connection.Close()

	client := services.NewPackagebirdServicesClient(connection)
	response, err := client.GetPackages(context.Background(), &services.Blank{})
	if err != nil {
		return err
	}

	for _, ele := range response.GetPackages() {
		fmt.Printf("%v-v%v\n", ele.GetName(), ele.GetVersion())
	}
	return nil
}

func getPackagesList() (*services.PackageList, error) {
	connection, err := grpc.Dial("127.0.0.1:55051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer connection.Close()

	client := services.NewPackagebirdServicesClient(connection)
	response, err := client.GetPackages(context.Background(), &services.Blank{})
	if err != nil {
		return nil, err
	}

	return response, nil
}
