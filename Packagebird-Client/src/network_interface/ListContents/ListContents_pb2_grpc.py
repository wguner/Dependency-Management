# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import src.network_interface.ListContents.ListContents_pb2 as ListContents__pb2


class ListContentServicesStub(object):
    """Serices associated with listing content operations
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetContent = channel.unary_unary(
                '/listcontent.ListContentServices/GetContent',
                request_serializer=ListContents__pb2.ContentRequest.SerializeToString,
                response_deserializer=ListContents__pb2.ContentResponse.FromString,
                )


class ListContentServicesServicer(object):
    """Serices associated with listing content operations
    """

    def GetContent(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ListContentServicesServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetContent': grpc.unary_unary_rpc_method_handler(
                    servicer.GetContent,
                    request_deserializer=ListContents__pb2.ContentRequest.FromString,
                    response_serializer=ListContents__pb2.ContentResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'listcontent.ListContentServices', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class ListContentServices(object):
    """Serices associated with listing content operations
    """

    @staticmethod
    def GetContent(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/listcontent.ListContentServices/GetContent',
            ListContents__pb2.ContentRequest.SerializeToString,
            ListContents__pb2.ContentResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
