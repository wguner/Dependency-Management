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
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"packagebird/network/grpc/services"

	"path/filepath"
)

// packageCmd represents the package command
var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("package called")
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
	connection, err := grpc.Dial("127.0.0.1:55051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer connection.Close()

	// Creating client
	client := services.NewPackagebirdServicesClient(connection)

	// Creating request
	project := &services.DownloadRequest{Path: name}

	// Create RPC stub
	download, err := client.DownloadFile(context.Background(), project)
	if err != nil {
		return err
	}

	// Create temp file
	tarfile, err := ioutil.TempFile(os.TempDir(), "src-")
	if err != nil {
		return err
	}
	// defer tarfile.Close()
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

	// Notify server
	if err := download.CloseSend(); err != nil {
		return err
	}

	// Create local directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create subdirectories
	separator := fmt.Sprintf("%c", os.PathSeparator)
	os.Mkdir(dir+separator+name, fs.ModePerm) // Project subdirectory
	_, err = tarfile.Seek(0, 0)
	if err != nil {
		return err
	}
	path, err := filepath.Abs(tarfile.Name())
	if err != nil {
		return err
	}
	cmd := exec.Command("tar", "-xf", path) // syscall.Exec(binary, args, os.Environ())
	err = cmd.Run()
	if err != nil {
		return err
	}

	tarfile.Close()

	// End operation
	return nil
}
