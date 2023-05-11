package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/scazaresm/go-file-storage-service/db"
	"github.com/scazaresm/go-file-storage-service/util"
)

func UploadFile(c *gin.Context) {
	// Get the file from the request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a unique ID for the file
	id := uuid.New().String()

	// Get the container name from the request parameters
	containerName := c.PostForm("containerName")
	if containerName == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing 'containerName' parameter."})
		return
	}

	// Get the folder (optional) from the request parameters, and normalize it
	folder := c.PostForm("folder")
	folder = strings.ReplaceAll(folder, "\\", "/")
	folder = strings.TrimPrefix(folder, "/")
	folder = strings.TrimSuffix(folder, "/")

	// Get the file name from the request parameters or use the original file name
	fileName := c.PostForm("fileName")
	if fileName == "" {
		fileName = file.Filename
	} else if strings.ContainsAny(fileName, "/\\") {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid file name, it cannot contain forward nor backward slashes."})
		return
	}

	// Save the file to disk
	err = c.SaveUploadedFile(file, fmt.Sprintf("data/%s/%s", containerName, id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Persist the file metadata to MongoDB
	_, err = db.SaveFileMetadata(db.FileMetadata{
		ID:        id,
		Container: containerName,
		Folder:    folder,
		FileName:  fileName,
		Created:   util.GetCurrentTime(),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
