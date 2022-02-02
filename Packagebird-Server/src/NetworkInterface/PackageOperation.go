package NetworkInterface

import (
	"context"
	databaseInterface "packagebird-server/src/DatabaseInterface"
	packageOperations "packagebird-server/src/NetworkInterface/PackageOperations"
	"packagebird-server/src/structures"
	"strconv"
)

func (server *GRPCServer) GetPackageList(context context.Context, request *packageOperations.PackageRequest) (*packageOperations.PackageListResponse, error) {
	// formatted as '$NAME-v$VERSION'
	// Really hacky solution. This is very poor and needs refactored
	pstring := databaseInterface.GetPackageNameVersion(request.GetPackageitem())
	requestedPackage := structures.Package{}
	requestedPackage.Name = pstring[0]
	requestedPackage.Version, _ = strconv.ParseInt(pstring[1], 10, 64)

	var packageList []string
	databaseInterface.GetPackageDependenciesRecurse(*mongoDBClientGlobal, requestedPackage.Name, int(requestedPackage.Version), &packageList)

	packageListResponseContents := &packageOperations.PackageListResponse{Packageitem: packageList}
	// packageListResponseContents.Packageitem = append(packageListResponseContents.Packageitem, packageList)

	return packageListResponseContents, nil
}
