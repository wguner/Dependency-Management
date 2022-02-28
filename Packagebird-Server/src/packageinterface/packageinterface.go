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
	dbinterface "packagebird-server/src/DatabaseInterface"
	networkinterface "packagebird-server/src/NetworkInterface"
	structures "packagebird-server/src/structures"
	"runtime"
	"strconv"
	"strings"

	extract "github.com/codeclysm/extract/v3"
)

func PackagePath(name string) string {
	return fmt.Sprintf("%v/%v.tar.gz", networkinterface.PACKAGEPATH, name)
}

func ResultPath(name string) string {
	return fmt.Sprintf("%v/%v/%v.exe", networkinterface.PACKAGEPATH, name, name)
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
	//pkg = pkgDataInter.GetPackage(*client, packageName, )
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

// Assuming that passed package is valid and complete, operation proceeds
func ExtractPackageSource(pkg *structures.Package, path string) (bool, error) {
	// Check if source file referenced exist
	if _, err := os.Stat(pkg.SourceFile); os.IsNotExist(err) {
		log.Printf("Cannot find source for package: %v-%d\nError:\t%v", pkg.Name, pkg.Version, err)
		return false, err
	}

	// Generate package extraction path
	pkgTmpDir := fmt.Sprintf("%v/%v-%v", path, pkg.Name, pkg.Version)

	// Check if subdirectory already present in tmp/, remove, return error if problem
	if _, err := os.Stat(pkgTmpDir); !os.IsNotExist(err) {
		if err := os.RemoveAll(pkgTmpDir); err != nil {
			log.Printf("Cannot remove: %v's directory and contents\nError:\t%v", pkgTmpDir, err)
			return false, err
		}
	}

	// Log to console
	log.Printf("Extracting package: %v", pkg.Name)

	// Create package subdirectory
	if err := os.Mkdir(pkgTmpDir, 0755); err != nil {
		log.Printf("Cannot create %v's directory\nError:\t%v", pkgTmpDir, err)
		return false, err
	}

	// Get File reference
	fd, err := ioutil.ReadFile(pkg.SourceFile)
	if err != nil {
		log.Printf("Cannot read file: %v\nError:\t%v", pkg.SourceFile, err)
		return false, err
	}

	// Extract archived source to package extraction directory
	err = extract.Gz(context.TODO(), bytes.NewBuffer(fd), pkgTmpDir, nil)
	if err != nil {
		log.Printf("Cannot extract file: %v\nError:\t%v", pkg.SourceFile, err)
		return false, err
	}

	// Successful operation, close
	return true, nil
}

// Executes the associated build process in the extracted package source directory
func BuildPackage(pkg *structures.Package, path string) (bool, error) {
	// Extracts, runs build commands
	if _, err := ExecutePackageCommand(pkg, path, pkg.BuildFile); err != nil {
		return false, err
	}
	return true, nil
}

// Executes associated test process in extracted package source directory
func TestPackage(pkg *structures.Package, path string) (bool, error) {
	// Extracts, runs test commands
	if _, err := ExecutePackageCommand(pkg, path, pkg.TestFile); err != nil {
		return false, err
	}
	return true, nil
}

// Executes a command on a package
func ExecutePackageCommand(pkg *structures.Package, path string, commands string) (bool, error) {
	// Extract package source
	if _, err := ExtractPackageSource(pkg, path); err != nil {
		return false, err
	}

	// Get Operating System to determine which shell
	var shell string
	os := runtime.GOOS
	switch os {
	case "windows":
		shell = "ps"
	case "linux":
		shell = "bash"
	default:
		shell = "bin/sh"
	}

	// Load commands to shell call in package directory
	cmd := exec.Command(shell, "-c", commands)
	cmd.Dir = fmt.Sprintf("%v/%v-%v", path, pkg.Name, pkg.Version)

	// Output for clarity
	log.Printf("Running: '%v'\nWith shell: '%v'\nOn directory: '%v'", commands, shell, cmd.Dir)

	// Run the command, capture error if present
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to run build command:\t'%v'\nOn package:\t'%v-%v'\nError:\t%v", commands, pkg.Name, pkg.Version, err)
		return false, err
	}

	// Output for success
	log.Printf("Successfully executed build associated with package.")

	return true, nil
}

// Extracts a package and all dependencies to a ./tmp/ subdirectory
func ExtractPackageDependenciesSource(pkg *structures.Package) error {
	pkgPath := fmt.Sprintf("./tmp/%v-%v", pkg.Name, pkg.Version)
	pkgDepsPath := fmt.Sprintf("%v/packages", pkgPath)

	// Extract base-level package
	if _, err := ExtractPackageSource(pkg, pkgPath); err != nil {
		return err
	}

	// Extract dependency-level packages
	log.Printf("Extracting package: %v dependencies", pkg.Name)
	for _, dep := range pkg.Dependencies {

		// Get dependency name, version
		depPkgName := strings.Split(dep, "-")[0]
		depPkgVer, err := strconv.Atoi(strings.Split(dep, "-")[1])
		if err != nil {
			log.Printf("Couldn't interpret string to number\nError:\t%v", err)
		}

		// Get dependency
		depPkg, err := dbinterface.GetPackage(*structures.GlobalMongoClient, depPkgName, depPkgVer)
		if err != nil {
			return err
		}

		// Perform extraction in ./tmp/pkg/ subdirectory
		depPkgExt := fmt.Sprintf("%v/%v", pkgDepsPath, depPkg.Name)
		if _, err := ExtractPackageSource(depPkg, depPkgExt); err != nil {
			return err
		}
	}
	return nil
}
