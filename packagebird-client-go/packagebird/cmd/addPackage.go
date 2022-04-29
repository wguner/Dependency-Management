/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"packagebird/network/grpc/services"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// addPackageCmd represents the addPackage command
var addPackageCmd = &cobra.Command{
	Use:   "package",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("add package called")
		pkgName := args[0]
		pkgVersion, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Print(err)
			return
		}

		err = addPackage(pkgName, pkgVersion)
		if err != nil {
			log.Print(err)
			return
		}

	},
}

func init() {
	addCmd.AddCommand(addPackageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addPackageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addPackageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addPackage(name string, version int64) error {
	// check if package exist
	exists, err := checkIfPackageInRegistry(name, version)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("package '%v-v%d' does not exist in registry", name, version)
	}

	// if successful, alter metadata and stats
	if err := alterPackageMetadata(name, version); err != nil {
		return err
	}

	// if successful, alter local file system
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	separator := fmt.Sprintf("%c", os.PathSeparator)
	proj := strings.TrimPrefix(wd, filepath.Dir(wd)+separator)
	if err := addToProject(proj, name, version); err != nil {
		return err
	}

	return nil
}

func checkIfPackageInRegistry(name string, version int64) (bool, error) {
	packages, err := getPackagesList()
	if err != nil {
		return false, err
	}
	passed := fmt.Sprintf("%v-v%d", name, version)

	for _, pck := range packages.GetPackages() {
		pckVal := fmt.Sprintf("%v-v%d", pck.GetName(), pck.GetVersion())
		if passed == pckVal {
			return true, nil
		}
	}
	return false, nil
}

func alterPackageMetadata(name string, version int64) error {
	client, connection, err := GetServerClient()
	defer connection.Close()
	_, err = client.UpdatePackageMetadata(context.Background(), &services.PackageRequest{Name: name, Version: version})
	if err != nil {
		return err
	}
	connection.Close()
	return nil
}

func addToProject(projectName string, packageName string, packageVersion int64) error {
	client, connection, err := GetServerClient()
	if err != nil {
		return err
	}

	r, err := client.AddPackage(context.Background(), &services.AddPackageRequest{
		ProjectName:    projectName,
		PackageName:    packageName,
		PackageVersion: packageVersion,
	})
	if err != nil {
		return err
	}

	// separator := fmt.Sprintf("%c", os.PathSeparator)
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	packagesDir := fmt.Sprintf(wd)
	packagesDir = filepath.FromSlash(fmt.Sprintf("%v%cpackages", packagesDir, os.PathSeparator))
	connection.Close()
	for _, path := range r.GetPaths() {
		path = filepath.FromSlash(path)
		if err := installPackage(path, packagesDir); err != nil {
			return err
		}
	}
	return nil
}

func installPackage(path string, loc string) error {
	// Create download connection to server
	client, connection, err := GetServerClient()
	defer connection.Close()
	if err != nil {
		return err
	}

	// Setup download operation
	download, err := client.DownloadData(context.Background(), &services.DownloadRequest{Path: path})
	if err != nil {
		return err
	}

	// Create temp file
	tarfile, err := ioutil.TempFile(os.TempDir(), "src-")
	if err != nil {
		return err
	}
	defer func(name string) error {
		err := os.Remove(name)
		if err != nil {
			return err
		}
		return nil
	}(tarfile.Name())

	// Begin receiving data
	for {
		// Receive chunk of download data
		chunk, err := download.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// End of download reached
		bytes := chunk.GetChunk()
		if (bytes == nil) || (len(bytes) == 0) {
			break
		}

		// Write data to temporary file
		_, err = tarfile.Write(bytes)
		if err != nil && err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}

	// Notify server of completion
	if err := download.CloseSend(); err != nil {
		return err
	}

	// Decompress, unwrap, install at location
	_, err = tarfile.Seek(0, 0) // Reset cursor to file beginning
	if err != nil {
		return err
	}
	tarPath, err := filepath.Abs(tarfile.Name())
	if err != nil {
		return err
	}
	cmd := exec.Command("tar", "-xf", tarPath, "-C", loc)
	if err = cmd.Run(); err != nil {
		return err
	}
	tarfile.Close()

	return nil
}
