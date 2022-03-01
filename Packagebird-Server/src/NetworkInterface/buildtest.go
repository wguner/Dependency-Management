package NetworkInterface

import (
	"fmt"
	"context"
	"log"
	buildtest "packagebird-server/src/NetworkInterface/buildtest"
	packageinterface "packagebird-server/src/packageinterface"
	structures "packagebird-server/src/structures"
	dbinterface "packagebird-server/src/DatabaseInterface"
)

func (server *GRPCServer) Build(context context.Context, pkgInfo *buildtest.PackageInfo) (*buildtest.BuildInfo, error) {
	pseudoPackage, err := dbinterface.GetPackage(*structures.GlobalMongoClient, pkgInfo.GetName(), int(pkgInfo.GetVersion()))
	if err != nil {
		log.Printf("Could not retreive package: %v", pkgInfo.GetName())
		return nil, err
	}

	// pseudoPackage := *structures.Package{}
	pseudoPackage.Name = pkgInfo.GetName()
	pseudoPackage.Version = pkgInfo.GetVersion()
	var buildInfoStatus string
	var buildInfoMessage string

	didBuild, err := packageinterface.BuildPackage(pseudoPackage);
	if err != nil {
		log.Printf("Err building package: %v-%v", pseudoPackage.Name, pseudoPackage.Version)
		return nil, err
	}

	if didBuild {
		buildInfoStatus = "OK"
		buildInfoMessage = fmt.Sprintf("Successfully built package: %v-%v", pseudoPackage.Name, pseudoPackage.Version)
		return &buildtest.BuildInfo{Status: buildInfoStatus, Message: buildInfoMessage}, nil
	} else {
		buildInfoStatus = "FAIL"
		buildInfoMessage = fmt.Sprintf("Failed to build package: %v-%v", pseudoPackage.Name, pseudoPackage.Version)
		return &buildtest.BuildInfo{Status: buildInfoStatus, Message: buildInfoMessage}, nil
	}
}

func (server *GRPCServer) Test(context context.Context, pkgInfo *buildtest.PackageInfo) (*buildtest.BuildInfo, error) {
	pseudoPackage, err := dbinterface.GetPackage(*structures.GlobalMongoClient, pkgInfo.GetName(), int(pkgInfo.GetVersion()))
	if err != nil {
		log.Printf("Could not retreive package: %v", pkgInfo.GetName())
		return nil, err
	}

	pseudoPackage.Name = pkgInfo.GetName()
	pseudoPackage.Version = pkgInfo.GetVersion()
	var buildInfoStatus string
	var buildInfoMessage string

	didBuild, err := packageinterface.TestPackage(pseudoPackage);
	if err != nil {
		log.Printf("Err testing package: %v-%v", pseudoPackage.Name, pseudoPackage.Version)
		return nil, err
	}

	if didBuild {
		buildInfoStatus = "OK"
		buildInfoMessage = fmt.Sprintf("Successfully tested package: %v-%v", pseudoPackage.Name, pseudoPackage.Version)
		return &buildtest.BuildInfo{Status: buildInfoStatus, Message: buildInfoMessage}, nil
	} else {
		buildInfoStatus = "FAIL"
		buildInfoMessage = fmt.Sprintf("Failed to test package: %v-%v", pseudoPackage.Name, pseudoPackage.Version)
		return &buildtest.BuildInfo{Status: buildInfoStatus, Message: buildInfoMessage}, nil
	}
}
