package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qunova/protein-analyzer/internal/models"
	"github.com/qunova/protein-analyzer/internal/services"
)

type UploadHandler struct {
	storageService *services.StorageService
}

func NewUploadHandler(storageService *services.StorageService) *UploadHandler {
	return &UploadHandler{
		storageService: storageService,
	}
}

// UploadFile handles protein structure file upload
func (h *UploadHandler) UploadFile(c *gin.Context) {
	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get file from request",
		})
		return
	}
	defer file.Close()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	validExtensions := map[string]bool{
		".pdb":   true,
		".cif":   true,
		".mmcif": true,
	}

	if !validExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file type. Only PDB, CIF, and mmCIF files are allowed",
		})
		return
	}

	// Validate file size (max 100MB)
	if header.Size > 100*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File size exceeds maximum limit of 100MB",
		})
		return
	}

	// Generate unique file ID
	fileID := uuid.New().String()

	// Save file
	savedPath, err := h.storageService.SaveFile(file, fileID, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to save file: %v", err),
		})
		return
	}

	// Prepare response
	response := models.UploadResponse{
		FileID:    fileID,
		FileName:  header.Filename,
		FileSize:  header.Size,
		FileType:  ext[1:], // Remove the dot
		UploadedAt: h.storageService.GetFileUploadTime(fileID),
	}

	c.JSON(http.StatusOK, response)
}