package packageinterface

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	extract "github.com/codeclysm/extract/v3"
)

func UnbundlePackage(packageName string) error {

	compressedFile, err := ioutil.ReadFile(packageName)
	if err != nil {
		log.Printf("Error reading package file for extraction")
		return err
	}

	dataFile := bytes.NewBuffer(compressedFile)
	newPackageName := strings.Split(packageName, "-")[0]
	extract.Gz(context.TODO(), dataFile, newPackageName, nil)

	return nil
}

func BuildAndCompressPackage(packageName string, language string) error {
	if err := RunBuildCommand(packageName, language); err != nil {
		return err
	}

	if err := CompressFile(packageName); err != nil {
		return err
	}

	return nil
}

func RunBuildCommand(packageName string, language string) error {
	if !FileExist(packageName) {
		return os.ErrNotExist
	}

	if language == "c" || language == "C" {
		commandString := fmt.Sprintf("gcc %v.c -o %v", packageName, packageName)
		err := exec.Command(commandString)
		if err != nil {
			log.Printf("Failed to build package name: %v", packageName)
			return nil
		}
	} else {
		log.Printf("Language %v not yet supported", language)
		return nil
	}
	return nil
}

func FileExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return true
	} else {
		log.Printf("Package %v not found at location", path)
		return false
	}
}

func CompressFile(fileName string) error {
	if !FileExist(fileName) {
		return os.ErrNotExist
	}

	fileRef, err := os.Open(fileName)
	if err != nil {
		log.Printf("Error getting file for compression")
		return err
	}

	fileRead := bufio.NewReader(fileRef)
	fileData, err := ioutil.ReadAll(fileRead)
	if err != nil {
		log.Printf("Error reading contents of file for compression")
		return err
	}

	// Bad, hardcodes for certain kind of executable
	compressedFileName := strings.Replace(fileName, ".exe", ".gz", -1)
	compressedFile, err := os.Create(compressedFileName)
	if err != nil {
		log.Printf("Error creating gzipped file for compression")
		return err
	}

	gzipWriter := gzip.NewWriter(compressedFile)
	gzipWriter.Write(fileData)
	gzipWriter.Close()

	return nil
}
