/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"packagebird/network/grpc/services"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "packagebird",
	Short: "A project and package management tool",
	Long: `Packagebird can create projects and packages synchronized with a remote registry setup by your System Admin.
You'll likely want to start with 'packagebird create project' followed by a project name, or 'packagebird create package'
if you're ready to create a package from the current iteration of a project.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.packagebird.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var GlobalServerAddress string

func GetServerClient() (services.PackagebirdServicesClient, *grpc.ClientConn, error) {
	connection, err := grpc.Dial(GetServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, err
	}
	client := services.NewPackagebirdServicesClient(connection)
	return client, connection, nil
}

func GetServerAddress() string {
	addr := os.Getenv("PACKAGEBIRD_SERVER_ADDRESS")
	if addr == "" || len(addr) <= 0 {
		os.Setenv("PACKAGEBIRD_SERVER_ADDRESS", "127.0.0.1:55051")
		return os.Getenv("PACKAGEBIRD_SERVER_ADDRESS")
	}
	return addr
}

var Client services.PackagebirdServicesClient
