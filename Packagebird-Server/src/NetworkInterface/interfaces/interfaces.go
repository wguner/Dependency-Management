package interfaces

import (
	fileTransfer "packagebird-server/src/NetworkInterface/FileTransfer"
	serverUtils "packagebird-server/src/NetworkInterface/ServerUtils"
	"packagebird-server/src/NetworkInterface/buildtest"
	"packagebird-server/src/NetworkInterface/listcontent"
	"packagebird-server/src/NetworkInterface/member_operations"
)

type GRPCServer struct {
	fileTransfer.UnimplementedFileServiceServer
	serverUtils.UnimplementedServerUtilsServicesServer
	listcontent.UnimplementedListContentServicesServer
	buildtest.UnimplementedBuildTestServicesServer
	member_operations.UnimplementedMemberCRUDServicesServer
	member_operations.UnimplementedMemberAuthenticationServer
}
