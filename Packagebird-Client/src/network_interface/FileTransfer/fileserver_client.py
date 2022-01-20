import grpc
import src.fileserver.FileTransfer_pb2 as FileTransfer_pb2
import src.fileserver.FileTransfer_pb2_grpc as FileTransfer_pb2_grpc

class FileTransfer(object):
    def __init__(self) -> None:
        pass

    # Parses content of file
    def file_reader(self, filename):
        BUFFER_SIZE = 64 * 1024
        first_chunk = True

        with open(filename, 'rb') as file:
            while True:
                if first_chunk:
                    # Minor change to ensure first 'chunk' submitted is the filename
                    yield FileTransfer_pb2.File(name=filename)
                    first_chunk = False

                chunks = file.read(BUFFER_SIZE)
                if len(chunks) == 0:
                    return
                yield FileTransfer_pb2.File(chunk=chunks)

    # Downloads file from server
    def download(self, address, port, filename):
        with grpc.insecure_channel(f'{address}:{port}') as channel:
            stub = FileTransfer_pb2_grpc.FileServiceStub(channel)
            response = stub.Download(FileTransfer_pb2.Request(body=f'{filename}'))
            with open(filename, 'wb') as file:
                for chunk in response:
                    file.write(chunk.chunk)
            print(f'Downloaded file: {filename}')

    # Uploads file to server
    def upload(self, address, port, filename):
        with grpc.insecure_channel(f'{address}:{port}') as channel:
            stub = FileTransfer_pb2_grpc.FileServiceStub(channel)
            file_chunks = self.file_reader(filename)
            response = stub.Upload(file_chunks)
            print(f'Uploaded file: {filename}, with response: {response.body}')
            nameFileRequest = FileTransfer_pb2.Request(header='Rename File', body=f'{filename}')
            response = stub.NameFile(nameFileRequest)
            print(f'Successfully renamed remote temp file to {filename},\nresponse: {response.body}')