import grpc
import src.network_interface.ServerUtils.ServerUtils_pb2 as ServerUtils_pb2
import src.network_interface.ServerUtils.ServerUtils_pb2_grpc as ServerUtils_pb2_grpc

class ServerUtils(object):
    
    # Pings the server to check for connection
    @staticmethod
    def ping(address, port) -> bool:
        with grpc.insecure_channel(f'{address}:{port}') as channel:
            stub = ServerUtils_pb2_grpc.ServerUtilsServicesStub(channel)
            clientinfo = ServerUtils_pb2.ClientInfo(body='Client request...')
            try:
                response = stub.Ping(clientinfo)
                if response.body != 'ACK':
                    print("Ping response indicates error...")
                    return False
                else:
                    return True
            except grpc.RpcError as rpcErr:
                if rpcErr.code() == grpc.StatusCode.UNAVAILABLE:
                    print("Couldn't connect to server; check client network configuration or insure that server is online...")
                else:
                    print("Unknown erorr encountered while pinging server...")
                return False