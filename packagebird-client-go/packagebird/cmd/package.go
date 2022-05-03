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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"packagebird/network/grpc/services"

	"path/filepath"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "project",
	Short: "Install a package from the server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := install(args[0]); err != nil {
			fmt.Print(err)
		}
	},
}

func init() {
	installCmd.AddCommand(packageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// packageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// packageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func install(name string) error {

	// Establishing connection
	connection, err := grpc.Dial(GetServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer connection.Close()

	// Creating client
	client := services.NewPackagebirdServicesClient(connection)

	// Creating request
	project := &services.DownloadRequest{Path: name}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := handleDownload(client, project, wd); err != nil {
		return err
	}

	// End operation
	return nil
}

func handleDownload(service services.PackagebirdServicesClient, request *services.DownloadRequest, path string) error {
	// Establish download handler
	download, err := service.DownloadFile(context.Background(), request)
	if err != nil {
		return err
	}

	// Create temporary file to receive data
	tarfile, err := ioutil.TempFile(os.TempDir(), "src-")
	if err != nil {
		return err
	}
	defer os.Remove(tarfile.Name())

	// Begin receiving data
	for {
		chunk, err := download.Recv()
		bytes := chunk.GetChunk()
		if (bytes == nil) || (len(bytes) == 0) {
			break
		}
		if err != nil {
			return err
		}

		_, err = tarfile.Write(bytes)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	// Notify server of finish
	err = download.CloseSend()
	if err != nil {
		return err
	}

	// Create or get directory specified in Path
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	// Get reference to temporary file
	tarfilePath, err := filepath.Abs(tarfile.Name())
	if err != nil {
		return err
	}

	// Execute extraction command, extract to specified location
	cmd := exec.Command("tar", "-xf", tarfilePath, "-C", path)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
