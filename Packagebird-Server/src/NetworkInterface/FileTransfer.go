package NetworkInterface

import (
	"fmt"
	"io"
	"log"
	"os"
	DatabaseInterface "packagebird-server/src/DatabaseInterface"
	fileTransfer "packagebird-server/src/NetworkInterface/FileTransfer"
	"packagebird-server/src/structures"
)

const (
	PACKAGEPATH = "C:\\Users\\ElishaAguilera\\Documents\\packages"
	CHUNKSIZE   = 64 * 1024
)

// Downloads a file from the server to the requesting client
func (server *GRPCServer) Download(request *fileTransfer.Request, fileStream fileTransfer.FileService_DownloadServer) error {
	log.Printf("Received request for package %v", request.GetBody())

	// Get direct path to file
	filepath := fmt.Sprintf("%vpackages\\%v", PACKAGEPATH, request.GetBody())

	// Open file and close when finished, or catch error
	file, err := os.Open(filepath)
	if err != nil {
		log.Printf("Error encountered opening file...")
		return err
	}
	defer file.Close()

	// Get chunk of bytes from file
	buffer := make([]byte, CHUNKSIZE)

	// Read all bytes in file up to chunk size
	for {
		// Read chunk of file
		bytes, err := file.Read(buffer)

		// Either end-of-file, or encountered error
		if err != nil {
			if err != io.EOF {
				log.Printf("Encountered error reading chunk of file...")
				return err
			} else {
				// Finished reading file, break out of loop
				break
			}
		}

		// Write chunk of file to struct
		filechunk := &fileTransfer.File{
			Content: &fileTransfer.File_Chunk{
				Chunk: buffer[:bytes],
			},
		}

		// Send file chunk to client, or error encountered
		err = fileStream.Send(filechunk)
		if err != nil {
			log.Printf("Encountered error transmitting file-chunk...")
			return err
		}
	}

	// Successfully completed file-download operation
	log.Printf("Succesfully downloaded file from server to client...")
	return nil
}

// Uploads a file from the client to the server
func (server *GRPCServer) Upload(fileStream fileTransfer.FileService_UploadServer) error {
	log.Printf("Received request to upload file from client to server...")

	// Creates the tempory directory if not already present
	var filepath = PACKAGEPATH // fmt.Sprintf("%vpackages", PACKAGEPATH)

	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		log.Printf("Temporary directory not found, creating new temporary directory")
		os.Mkdir(filepath, 0755)
	}

	// Get the name of the file from the first 'chunk' received
	chunk, err := fileStream.Recv()
	if err != nil {
		log.Printf("Error encountered receiving first chunk (filename) from client...")
		return err
	}
	filename := chunk.GetName()

	// Create the temp file to be written to
	file, err := os.Create(filepath + "\\" + filename)

	// If error encountered, return out. Else, progress to writing to the file
	if err != nil {
		log.Printf("Error encountered attempting to write to file...")
		return err
	}
	defer file.Close()

	// Iterate until chunks are no longer received from client
	for {
		chunk, err := fileStream.Recv()

		// If chunk is empty, end of file reached, break out of loop
		if (chunk == nil) || (len(chunk.GetChunk()) == 0) {
			break
		}

		// If error encountered, exit operation
		if err != nil {
			log.Printf("Error encountered receiving file from client...")
			return err
		}

		// Otherwise, write chunk to file
		_, err = file.Write(chunk.GetChunk())
		if err != nil {
			// Either end-of-file reached, or another error encountered
			if err == io.EOF {
				break
			}
			log.Printf("Error encountered writing file chunk to file...")
			return err
		}
	}

	entry := structures.Package{
		Name:    filename,
		Version: 0,
	}
	_, err = DatabaseInterface.NewPackage(*mongoDBClientGlobal, entry)

	var responseFileName string
	var message *fileTransfer.Response

	// Respond to client
	if err != nil {
		responseFileName = fmt.Sprintf("File %v uploaded successfully...", filename)
	} else {
		responseFileName = fmt.Sprintf("File %v did not upload successfully due to database error...", filename)
	}

	// Respond to client informing of success
	message = &fileTransfer.Response{
		Body: responseFileName,
	}

	// Inform client of success, close operation
	fileStream.SendAndClose(message)

	log.Printf("File upload operation completed...")
	return nil
}
