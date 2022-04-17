/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Creates a new packagebird project, with code contributions committed therein",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating new project directory")
		createProject(args)
	},
}

func init() {
	createCmd.AddCommand(projectCmd)
}

func createProject(projectName []string) {
	if err := os.Mkdir(projectName[0], fs.ModePerm); err != nil {
		log.Fatal("Error creating project directory")
	}
}
