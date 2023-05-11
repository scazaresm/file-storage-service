package client

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func createFileHeader(filePath string) (*multipart.FileHeader, error) {
	// Open the file and get its information
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// Create a new file header with the file name and size
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", filepath.Base(filePath)))
	header.Set("Content-Type", "application/octet-stream")
	header.Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))

	return &multipart.FileHeader{
		Filename: filepath.Base(filePath),
		Header:   header,
		Size:     fileInfo.Size(),
	}, nil
}

func TestUploadFile(t *testing.T) {
	file, err := createFileHeader("../README.md")
	if err != nil {
		log.Fatal(err)
	}

	req := UploadFileRequest{
		File:          file,
		FileName:      "my-file.txt",
		ContainerName: "my-container",
		Folder:        "test/foo/bar",
	}

	c := NewFileStorageClient(&http.Client{}, "http://localhost:8081")
	res, err := c.UploadFile(req)

	b, err := io.ReadAll(res.Body)
	fmt.Println(string(b))

}
