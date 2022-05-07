/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"packagebird/network/grpc/services"
	"strconv"

	"github.com/spf13/cobra"
)

// removePackageCmd represents the removePackage command
var removePackageCmd = &cobra.Command{
	Use:   "remove",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("removePackage called")
		version, _ := strconv.ParseInt(args[2], 10, 64)
		removePackage(args[0], args[1], version)
	},
}

func init() {
	removeCmd.AddCommand(removePackageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removePackageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removePackageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func removePackage(projectName string, packageName string, packageVersion int64) error {
	client, connection, err := GetServerClient()
	defer connection.Close()
	if err != nil {
		return err
	}
	r, err := client.RemovePackage(context.Background(), &services.AddPackageRequest{
		ProjectName:    projectName,
		PackageName:    packageName,
		PackageVersion: packageVersion,
	})
	if err != nil {
		return err
	}
	if !r.GetSuccess() {
		return nil
	}
	wd, _ := os.Getwd()
	path := fmt.Sprintf("%v%c%v-v%d", wd, os.PathSeparator, packageName, packageVersion)
	if err := os.Remove(path); err != nil {
		return err
	}
	return nil
}
