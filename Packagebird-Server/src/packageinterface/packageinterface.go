package packageinterface

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	networkinterface "packagebird-server/src/NetworkInterface"
	"strings"

	extract "github.com/codeclysm/extract/v3"
)

func PackagePath(name string) string {
	return fmt.Sprintf("%v\\%v.tar.gz", networkinterface.PACKAGEPATH, name)
}

func ResultPath(name string) string {
	return fmt.Sprintf("%v\\%v\\%v.exe", networkinterface.PACKAGEPATH, name, name)
}

func UnbundlePackage(packageName string) error {
	compressedFile, err := ioutil.ReadFile(PackagePath(packageName))
	if err != nil {
		log.Printf("Error reading package file for extraction")
		return err
	}

	dataFile := bytes.NewBuffer(compressedFile)
	newPackageName := strings.Split(packageName, "-")[0]

	err = extract.Gz(context.TODO(), dataFile, newPackageName, nil)
	if err != nil {
		return err
	}

	return nil
}

func BuildAndCompressPackage(packageName string, language string) error {
	if err := RunBuildCommand(BuildPath(packageName), language); err != nil {
		return err
	}

	if err := CompressFile(BuildPath(packageName)); err != nil {
		return err
	}

	return nil
}

func BuildPath(name string) string {
	return fmt.Sprintf(`%v\%v\%v`, networkinterface.PACKAGEPATH, name, name)
}

func RunBuildCommand(packageName string, language string) error {
	/*if !FileExist(packageName) {
		return os.ErrNotExist
	}*/

	if language == "c" || language == "C" {
		commandString := fmt.Sprintf(`gcc %v.c -o %v`, BuildPath(packageName), BuildPath(packageName))
		cmd := exec.Command(commandString)
		err := cmd.Run()
		if err != nil {
			log.Printf("Failed to build package name: %v", packageName)
			return err
		}
	} else {
		log.Printf("Language %v not yet supported", language)
		return errors.New("cannot build package")
	}
	return nil
}

func FileExist(name string) bool {
	if _, err := os.Stat(PackagePath(name)); err != nil {
		return true
	} else {
		log.Printf("Package %v not found at location", name)
		return false
	}
}

func CompressFile(fileName string) error {
	/*if !FileExist(fileName) {
		return os.ErrNotExist
	}*/

	fileRef, err := os.Open(ResultPath(fileName))
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
	compressedFileName := strings.Replace(ResultPath(fileName), ".exe", ".gz", -1)
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
