package internal

import (
	"io"
	"os"
)

func WriteToDisk(respBody io.Reader, outFile *os.File) error {
	const bufferSize = 4096
	buffer := make([]byte, bufferSize)
	for {
		bytesRead, err := respBody.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.EOF {
			return err
		}
		if bytesRead == 0 {
			break
		}
		if _, err := outFile.Write(buffer[:bytesRead]); err != nil {
			return err
		}
	}
	return nil
}

func WriteBytesToDisk(b []byte, outFile *os.File) error {
	_, err := outFile.Write(b)
	return err
}
