from turtle import mode
import grpc
import src.network_interface.FileTransfer.FileTransfer_pb2 as FileTransfer_pb2
import src.network_interface.FileTransfer.FileTransfer_pb2_grpc as FileTransfer_pb2_grpc
import datetime

class FileTransfer(object):
    def __init__(self) -> None:
        pass

    # Parses content of file
    def file_reader(self, filename, operationMode):
        BUFFER_SIZE = 64 * 1024
        mode_need_shared = True
        first_chunk = True

        print(f'File name being uploaded: {filename}\nMode of upload: {operationMode}\n')
        start_time = datetime.datetime.now()
        
        with open(filename, 'rb') as file:
            while True:
                # print(f'Values of control booleans:\nMode:\t{mode_need_shared}\nFirst Chunk:\t{first_chunk}')
                if mode_need_shared:
                    print(f'Output upload type {operationMode}')
                    yield FileTransfer_pb2.File(type=operationMode)
                    mode_need_shared=False
                elif first_chunk:
                    # Minor change to ensure first 'chunk' submitted is the filename
                    yield FileTransfer_pb2.File(name=filename)
                    first_chunk = False
                else:
                    chunks = file.read(BUFFER_SIZE)
                    if len(chunks) == 0:
                        return
                    delta_time = datetime.datetime.now() - start_time
                    total_seconds = delta_time.total_seconds()
                    print(f"Total seconds since upload start: {total_seconds}", end='\r')
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
    def upload(self, address, port, filename, mode):
        with grpc.insecure_channel(f'{address}:{port}') as channel:
            stub = FileTransfer_pb2_grpc.FileServiceStub(channel)
            file_chunks = self.file_reader(filename, mode)
            # print(f'Chunk being uploaded: {file_chunks}')
            response = stub.Upload(file_chunks)
            # print(f'Uploaded file: {filename}, with response: {response.body}')
            # nameFileRequest = FileTransfer_pb2.Request(header='Rename File', body=f'{filename}')
            # response = stub.NameFile(nameFileRequest)
            # print(f'Successfully renamed remote temp file to {filename},\nresponse: {response.body}')