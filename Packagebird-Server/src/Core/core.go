package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"packagebird-server/src/DatabaseInterface"
	NetworkInterface "packagebird-server/src/networkinterface"

	yaml "gopkg.in/yaml.v3"
)

// Temporary constant for address
const (
	ADDRESS  = "127.0.0.1:50051"
	MONGOURI = "mongodb://localhost:27017"
)

// Global configuration structure
var config *ConfigFile

// Entry-poin for the server application
func main() {
	// Loads from configuration file
	config, err := setupConfigFile("config.yaml")
	if err != nil {
		log.Print(err)
	}

	mongoDBClient, err := DatabaseInterface.MongoDBServerConnect(config.DatabaseAddress)
	if err != nil {
		log.Fatalf("Error encountered on server connecting to MongoDB:\n%v\nShutting down...", err)
	}

	if err := NetworkInterface.PackagebirdServerStart(config.ServerAddress, mongoDBClient); err != nil {
		log.Fatalf("Error encountered on gRPC server:\n%v\nShutting down...", err)
	}

	err = DatabaseInterface.MongoDBServerDisconnect(*mongoDBClient)
	if err != nil {
		log.Fatalf("Error encountered on server disconnecting from MongoDB:\n%v\nShutting down...", err)
	}
}

type ConfigFile struct {
	PackageSourcePath string
	ProjectSourcePath string
	ServerAddress     string
	DatabaseAddress   string
}

func setupConfigFile(path string) (*ConfigFile, error) {
	// Check if configuration file exist
	if _, err := os.Stat(path); err != nil {
		// If not, create default configuration file
		log.Printf("Config file: %v was not found in expected location.\nCreating default configuration file", path)
		err = createConfigFile(path)
		// Horrible way of error handling
		if err != nil {
			return nil, err
		}
	}

	// Load configuration file to global variable
	file, err := getConfigFile(path)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func createConfigFile(path string) error {
	// Get current working directory of server
	wd, err := os.Getwd()
	if err != nil {
		log.Printf("Cannot get current working directory path\nError: %v", err)
		return err
	}

	// Default packages, projects directory in same working directory as server
	file := ConfigFile{
		PackageSourcePath: fmt.Sprintf("%v\\packages", wd),
		ProjectSourcePath: fmt.Sprintf("%v\\projects", wd),
		ServerAddress:     "127.0.0.1:50051",
		DatabaseAddress:   "mongodb://localhost:27017",
	}

	// Encode structure data to YAML
	yaml, err := yaml.Marshal(&file)
	if err != nil {
		log.Printf("Cannot marshall default config file to YAML\nError: %v", err)
		return err
	}

	// Write to local config YAML file
	err = ioutil.WriteFile(path, yaml, 0644)
	if err != nil {
		log.Printf("Cannot write to local file\nError: %v", err)
		return err
	}

	return nil
}

func getConfigFile(path string) (*ConfigFile, error) {
	// Read contents of local file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read contents of file: %v\nError: %v", path, err)
		return nil, err
	}

	// Create structure mapping
	configFile := ConfigFile{}

	// Parse file as configuration file in YAML format
	err = yaml.Unmarshal(file, &configFile)
	if err != nil {
		log.Printf("Failed to map contents of file: %v to ConfigFile structure\nError: %v", path, err)
		return nil, err
	}

	return &configFile, nil
}
