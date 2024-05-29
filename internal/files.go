package internal

import (
	"fmt"
	"io"
	"net/http"
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

func DownloadToDisk(url string, outFile *os.File) error {
	// 从URL获取数据
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while getting the URL: %v", err)
	}
	defer resp.Body.Close()

	// 将正文写入文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("error while writing to the file: %v", err)
	}

	return nil
}
