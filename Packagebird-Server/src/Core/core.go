package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"packagebird-server/src/DatabaseInterface"
	configFile "packagebird-server/src/config"
	NetworkInterface "packagebird-server/src/networkinterface"
	structures "packagebird-server/src/structures"

	yaml "gopkg.in/yaml.v3"
)

var GlobalConfig *configFile.ConfigFile

// Temporary constant for address
const (
	ADDRESS  = "127.0.0.1:50051"
	MONGOURI = "mongodb://localhost:27017"
)

// Entry-poin for the server application
func main() {
	// Loads from configuration file
	Config, err := setupConfigFile("config.yaml")
	if err != nil {
		log.Print(err)
	}

	// Set global package variable
	configFile.Config = Config

	// Setup Server subdirectories
	if err := createServerSubdirectories(); err != nil {
		log.Fatalf("Error encountered creating or finding required server subdirectories:\t%v\nShutting down...", err)
	}

	structures.GlobalMongoClient, err = DatabaseInterface.MongoDBServerConnect(Config.DatabaseAddress)
	if err != nil {
		log.Fatalf("Error encountered on server connecting to MongoDB:\n%v\nShutting down...", err)
	}

	if err := NetworkInterface.PackagebirdServerStart(Config.ServerAddress, structures.GlobalMongoClient); err != nil {
		log.Fatalf("Error encountered on gRPC server:\n%v\nShutting down...", err)
	}

	err = DatabaseInterface.MongoDBServerDisconnect(*structures.GlobalMongoClient)
	if err != nil {
		log.Fatalf("Error encountered on server disconnecting from MongoDB:\n%v\nShutting down...", err)
	}
}

func setupConfigFile(path string) (*configFile.ConfigFile, error) {
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
	file := configFile.ConfigFile{
		PackageSourcePath: fmt.Sprintf("%v/packages", wd),
		ProjectSourcePath: fmt.Sprintf("%v/projects", wd),
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

func getConfigFile(path string) (*configFile.ConfigFile, error) {
	// Read contents of local file
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read contents of file: %v\nError: %v", path, err)
		return nil, err
	}

	// Create structure mapping
	configFile := configFile.ConfigFile{}

	// Parse file as configuration file in YAML format
	err = yaml.Unmarshal(file, &configFile)
	if err != nil {
		log.Printf("Failed to map contents of file: %v to ConfigFile structure\nError: %v", path, err)
		return nil, err
	}

	return &configFile, nil
}

// Creates expected subdirectories for server use if not already present
func createServerSubdirectories() error {
	subdirectories := []string{"projects", "packages", "builds", "tmp"}
	if err := createSubdirectories(subdirectories); err != nil {
		return err
	}
	return nil
}

// Creates subdirectories within working directory if not already present
func createSubdirectories(paths []string) error {
	for _, path := range paths {
		if err := createSubdirectory(path); err != nil {
			return err
		}
	}
	return nil
}

// Creates subdirectory within working directory if not already present
func createSubdirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
		return nil
	} else if !os.IsNotExist(err) {
		return nil
	} else {
		log.Print(err)
		return err
	}
}
