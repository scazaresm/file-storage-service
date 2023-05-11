package client

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

type FileStorageClient struct {
	Client  *http.Client
	BaseUrl string
}

type UploadFileRequest struct {
	File          *multipart.FileHeader
	FileName      string
	ContainerName string
	Folder        string
}

func NewFileStorageClient(client *http.Client, baseURL string) *FileStorageClient {
	return &FileStorageClient{
		Client:  client,
		BaseUrl: baseURL,
	}
}

func (c FileStorageClient) UploadFile(req UploadFileRequest) (*http.Response, error) {
	// Open the file
	fileContent, err := req.File.Open()
	if err != nil {
		return nil, err
	}
	defer fileContent.Close()

	// Create a new multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the file to the form
	part, err := writer.CreateFormFile("file", req.File.Filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_, err = io.Copy(part, fileContent)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Add other fields to the form
	writer.WriteField("containerName", req.ContainerName)
	writer.WriteField("folder", req.Folder)
	writer.WriteField("fileName", req.FileName)

	err = writer.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	url := c.BaseUrl + "/v1/file"

	// Create a new HTTP request
	httpReq, err := http.NewRequest("POST", url, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	resp, err := c.Client.Do(httpReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return resp, nil
}
