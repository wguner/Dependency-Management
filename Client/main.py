import grpc
import FileTransfer_pb2
import FileTransfer_pb2_grpc

def file_reader(filename):
    BUFFER_SIZE = 64 * 1024
    
    # Workaround to write filename to first chunk
    '''
    tempFile = open(f'temp_{filename}', 'wb')
    filename_bytes = len(filename.encode('utf-8'))
    chunkOffset = BUFFER_SIZE - filename_bytes
    first_chunk = b""+f'{filename}'+bytes(chunkOffset)
    tempFile.write(first_chunk)

    with open(filename, 'rb') as file:
        while True:
            chunk = file.read(BUFFER_SIZE)
            if len(chunk) == 0:
                break
            tempFile.write(chunk)
    
    tempFile.close()    
    '''

    with open(filename, 'rb') as file:
        while True:
            chunks = file.read(BUFFER_SIZE)
            if len(chunks) == 0:
                return
            yield FileTransfer_pb2.File(chunk=chunks)

def download(filename):
    with grpc.insecure_channel('127.0.0.1:50051') as channel:
        stub = FileTransfer_pb2_grpc.FileServiceStub(channel)
        response = stub.Download(FileTransfer_pb2.Request(body=f'{filename}'))
        with open(filename, 'wb') as file:
            for chunk in response:
                file.write(chunk.chunk)
        print(f'Downloaded file: {filename}')

def upload(filename):
    with grpc.insecure_channel('127.0.0.1:50051') as channel:
        stub = FileTransfer_pb2_grpc.FileServiceStub(channel)
        file_chunks = file_reader(filename)
        # initial = stub.Upload(FileTransfer_pb2.File(name=filename))
        response = stub.Upload(file_chunks)
        print(f'Uploaded file: {filename}, with response: {response}')
        nameFileRequest = FileTransfer_pb2.Request(header='Rename File', body=f'{filename}')
        response = stub.NameFile(nameFileRequest)
        print(f'Successfully renamed remote temp file to {filename},\nresponse: {response}')

if __name__ == '__main__':
    # download('Test.txt')
    upload('test.txt')