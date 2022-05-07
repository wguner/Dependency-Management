/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

// getProjectPackagesCmd represents the getProjectPackages command
var getProjectPackagesCmd = &cobra.Command{
	Use:   "added",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := getInstalledPackages()
		if err != nil {
			log.Print(err)
		}
	},
}

func init() {
	getCmd.AddCommand(getProjectPackagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getProjectPackagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getProjectPackagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func getInstalledPackages() error {
	packagesFolder := fmt.Sprintf(".%cpackages", os.PathSeparator)
	_, err := os.Stat(packagesFolder)
	if err != nil {
		return err
	}
	dirs, err := ioutil.ReadDir(packagesFolder)
	if err != nil {
		return err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			fmt.Printf("%v\n", dir.Name())
		}
	}
	return nil
}
