/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"packagebird/network/grpc/services"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// createPackageCmd represents the createPackage command
var createPackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Creates a new package from a project",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := strconv.ParseInt(args[1], 10, 64)
		createPackage(args[0], version)
	},
}

func init() {
	createCmd.AddCommand(createPackageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createPackageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createPackageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createPackage(name string, version int64) error {

	// Create bundled directory
	// tarfile, err := makeDirectoryTarfile()
	/* if err != nil {
		return err
	} */

	// Establish connection
	connection, err := grpc.Dial(GetServerAddress(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer connection.Close()

	// Create client
	client := services.NewPackagebirdServicesClient(connection)

	// Create package template
	pkg, err := client.CreatePackage(context.Background(), &services.PackageRequest{
		Name:    name,
		Version: version,
	})

	/*
		// Create RPC stub
		upload, err := client.UploadFile(context.Background())
		if err != nil {
			return err
		}

		// Send Path
		path := pkg.GetMessage()
		file := &services.File{
			Content: &services.File_Path{
				Path: path,
			},
		}
		if err := upload.Send(file); err != nil {
			log.Print(err)
			return err
		}

		// Transmit package source
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
		tarfile.Close()
		response, err := upload.CloseAndRecv()
		if err != nil {
			fmt.Print(err)
			return err
		}
	*/

	fmt.Print(pkg)

	return nil
}

func makeDirectoryTarfile() (*os.File, error) {

	// Get working directory
	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Create virtual filesystem
	filesystem := os.DirFS(root)

	// Create temporary archive for transmission
	tarfile, err := ioutil.TempFile(os.TempDir(), "src-")
	if err != nil {
		return nil, err
	}

	// Walk project directory, get all subdirectories and files NOT in '/packages'
	dirs := make(map[string]string)
	if err := fs.WalkDir(filesystem, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() && path == "packages" {
			return fs.SkipDir
		}
		if !entry.IsDir() {
			abs, _ := filepath.Abs(entry.Name())
			separator := fmt.Sprintf("%c", os.PathSeparator)
			dirs[path] = strings.TrimPrefix(abs, filepath.Dir(root)+separator)
			if err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Make archive
	files, err := archiver.FilesFromDisk(nil, dirs)
	if err != nil {
		return nil, err
	}
	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}
	if err := format.Archive(context.Background(), tarfile, files); err != nil {
		log.Print(err)
		return nil, err
	}

	// Reset cursor to beginning
	_, err = tarfile.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	// Return generated archive
	return tarfile, nil
}
