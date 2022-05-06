package NetworkInterface

import (
	"context"
	"fmt"
	databaseInterface "packagebird-server/src/DatabaseInterface"
	packageOperations "packagebird-server/src/NetworkInterface/PackageOperations"
	structures "packagebird-server/src/structures"
	"strconv"
)

func (server *GRPCServer) GetPackageList(context context.Context, request *packageOperations.PackageRequest) (*packageOperations.PackageListResponse, error) {
	// formatted as '$PNAME-$NAME-v$VERSION'
	// Really hacky solution. This is very poor and needs refactored
	pstring := databaseInterface.GetPackageNameVersion(request.GetPackageitem())

	// Project name
	pr_name := pstring[0]

	// Lookup project
	pr_exist, err := databaseInterface.LookupProject(*mongoDBClientGlobal, pr_name, "")
	if !pr_exist {
		return nil, err
	}

	requestedPackage := structures.Package{}
	requestedPackage.Name = pstring[1]
	requestedPackage.Version, _ = strconv.ParseInt(pstring[2], 10, 64)

	var packageList []string
	databaseInterface.GetPackageDependenciesRecurse(*mongoDBClientGlobal, requestedPackage.Name, int(requestedPackage.Version), &packageList)

	// Should update project dependencies
	depstring := fmt.Sprintf("%v-v%v", pstring[1], pstring[2])
	deparr := []string{depstring}
	err = databaseInterface.UpdateProject(*mongoDBClientGlobal, pr_name, deparr)
	if err != nil {
		return nil, err
	}

	packageListResponseContents := &packageOperations.PackageListResponse{Packageitem: packageList}
	// packageListResponseContents.Packageitem = append(packageListResponseContents.Packageitem, packageList)

	return packageListResponseContents, nil
}
