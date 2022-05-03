/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// setServerCmd represents the setServer command
var setServerCmd = &cobra.Command{
	Use:   "setServer",
	Short: "Set the server attribute",
	Run: func(cmd *cobra.Command, args []string) {
		SetServer(args[0])
	},
}

func init() {
	setCmd.AddCommand(setServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func SetServer(address string) {
	os.Setenv("PACKAGEBIRD_SERVER_ADDRESS", address)
}
