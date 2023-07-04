package fs

import (
	"io"
	"os"
)

// chunkSize defines the size of the chunks to read
const chunkSize = 4096

// SaveToDisk uploads the contents of a multipart file to disk with the given fileName.
// If the file already exists in the file system, an error is returned.
// The function copies the contents of the multipart file to the new file in chunks
// of size chunkSize. It returns an error if any I/O error occurs during the upload.
func SaveToDisk(f io.ReadSeekCloser, fileName string) error {
	defer f.Close()
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		return err
	}
	fileInDisk, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer fileInDisk.Close()
	// copy the contents of the multipart.Bytes to the new file in chunks
	chunkSize := chunkSize
	buf := make([]byte, chunkSize)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		_, err = fileInDisk.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	return nil
}
