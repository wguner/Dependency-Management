/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	services "packagebird/network/grpc/services"
	"path/filepath"
	"strings"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs the project to the current version in the server.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("syncing project directory")

		wd, err := os.Getwd()
		dir := wd
		if err != nil {
			log.Print(err)
		}

		err = sync(dir)
		if err != nil {
			log.Print(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// syncCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// syncCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sync(name string) error {

	// Get project directory path
	root, err := os.Getwd()
	if err != nil {
		return err
	}

	// Create virtual file system
	filesystem := os.DirFS(root)
	if err != nil {
		return err
	}

	// Create temporary archive for transmission
	tarfile, err := ioutil.TempFile(os.TempDir(), "src-")
	if err != nil {
		return err
	}
	defer tarfile.Close()

	// Walk project directory, get all subdirectories and files NOT in '/packages'
	dirs := make(map[string]string)

	if err := fs.WalkDir(filesystem, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, "packages") {
			return fs.SkipDir
		}
		path = filepath.FromSlash(path)
		fmt.Println("Walking: ", path, "Dir?", entry.IsDir())
		if !entry.IsDir() {
			// abs, _ := filepath.Abs(entry.Name())
			separator := fmt.Sprintf("%c", os.PathSeparator)
			dirs[path] = strings.TrimPrefix(root+separator+path, filepath.Dir(root)+separator)
			if err != nil {
				return err
			}
			fmt.Println("Adding a file with path", path, "and name", entry.Name())
		}
		return nil
	}); err != nil {
		return err
	}

	// Make archive
	files, err := archiver.FilesFromDisk(nil, dirs)
	if err != nil {
		return err
	}
	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}
	if err := format.Archive(context.Background(), tarfile, files); err != nil {
		log.Print(err)
		return err
	}

	// Reset cursor to beginning
	_, err = tarfile.Seek(0, 0)
	if err != nil {
		return err
	}

	// Establishing connection
	connection, err := grpc.Dial("127.0.0.1:55051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer connection.Close()

	// Creating client
	client := services.NewPackagebirdServicesClient(connection)

	// Create RPC stub
	upload, err := client.UploadFile(context.Background())
	if err != nil {
		return err
	}

	// Send path
	separator := fmt.Sprintf("%c", os.PathSeparator)
	dir := strings.TrimPrefix(root, filepath.Dir(root)+separator)
	file := &services.File{Content: &services.File_Path{Path: dir}}
	if err := upload.Send(file); err != nil {
		log.Print(err)
		return err
	}

	// Transmitting file
	buffer := make([]byte, 64*1024)
	for {
		bytes, err := tarfile.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		chunk := &services.File{
			Content: &services.File_Chunk{
				Chunk: buffer[:bytes],
			},
		}

		err = upload.Send(chunk)
		if err != nil {
			return err
		}
	}

	// Getting response
	response, err := upload.CloseAndRecv()
	if err != nil {
		return err
	}
	if response.GetSuccess() {
		fmt.Printf(response.GetMessage())
	}

	return nil
}
