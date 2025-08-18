package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/qunova/protein-analyzer/internal/models"
	"github.com/qunova/protein-analyzer/internal/services"
)

type AnalysisHandler struct {
	analysisService *services.AnalysisService
	storageService  *services.StorageService
}

func NewAnalysisHandler(analysisService *services.AnalysisService, storageService *services.StorageService) *AnalysisHandler {
	return &AnalysisHandler{
		analysisService: analysisService,
		storageService:  storageService,
	}
}

// AnalyzeStructure initiates protein structure analysis
func (h *AnalysisHandler) AnalyzeStructure(c *gin.Context) {
	fileID := c.Param("fileId")
	
	var request models.AnalysisRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		// Use default options if not provided
		request = models.AnalysisRequest{
			FileID: fileID,
			Options: models.AnalysisOptions{
				IncludeSequence:  true,
				CalculateDomains: true,
				IdentifyLigands:  true,
				Resolution:       "high",
			},
		}
	}

	// Check if file exists
	filePath, err := h.storageService.GetFilePath(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "File not found",
		})
		return
	}

	// Generate analysis ID
	analysisID := uuid.New().String()

	// Start analysis asynchronously
	go h.analysisService.ProcessStructure(analysisID, filePath, request.Options)

	// Return immediate response
	c.JSON(http.StatusAccepted, gin.H{
		"analysisId": analysisID,
		"status":     "processing",
		"message":    "Analysis started successfully",
	})
}

// GetAnalysisResult retrieves the analysis result
func (h *AnalysisHandler) GetAnalysisResult(c *gin.Context) {
	analysisID := c.Param("analysisId")

	result, err := h.analysisService.GetAnalysisResult(analysisID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Analysis not found",
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetChains retrieves chain data from analysis
func (h *AnalysisHandler) GetChains(c *gin.Context) {
	analysisID := c.Param("analysisId")

	result, err := h.analysisService.GetAnalysisResult(analysisID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Analysis not found",
		})
		return
	}

	if result.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Analysis not completed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"analysisId": analysisID,
		"chains":     result.Structure.Chains,
		"count":      len(result.Structure.Chains),
	})
}

// GetDomains retrieves domain data from analysis
func (h *AnalysisHandler) GetDomains(c *gin.Context) {
	analysisID := c.Param("analysisId")

	result, err := h.analysisService.GetAnalysisResult(analysisID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Analysis not found",
		})
		return
	}

	if result.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Analysis not completed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"analysisId": analysisID,
		"domains":    result.Structure.Domains,
		"count":      len(result.Structure.Domains),
	})
}

// GetLigands retrieves ligand data from analysis
func (h *AnalysisHandler) GetLigands(c *gin.Context) {
	analysisID := c.Param("analysisId")

	result, err := h.analysisService.GetAnalysisResult(analysisID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Analysis not found",
		})
		return
	}

	if result.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Analysis not completed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"analysisId": analysisID,
		"ligands":    result.Structure.Ligands,
		"count":      len(result.Structure.Ligands),
	})
}

// ExportData exports all analysis data
func (h *AnalysisHandler) ExportData(c *gin.Context) {
	analysisID := c.Param("analysisId")

	result, err := h.analysisService.GetAnalysisResult(analysisID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Analysis not found",
		})
		return
	}

	if result.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Analysis not completed",
		})
		return
	}

	// Set headers for file download
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"analysis_%s.json\"", analysisID))

	c.JSON(http.StatusOK, result.Structure)
}

// ExportSpecificData exports specific type of data (chains, domains, or ligands)
func (h *AnalysisHandler) ExportSpecificData(c *gin.Context) {
	analysisID := c.Param("analysisId")
	dataType := c.Param("dataType")

	result, err := h.analysisService.GetAnalysisResult(analysisID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Analysis not found",
		})
		return
	}

	if result.Status != "completed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Analysis not completed",
		})
		return
	}

	var data interface{}
	switch dataType {
	case "chains":
		data = result.Structure.Chains
	case "domains":
		data = result.Structure.Domains
	case "ligands":
		data = result.Structure.Ligands
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid data type. Use 'chains', 'domains', or 'ligands'",
		})
		return
	}

	// Set headers for file download
	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s_%s.json\"", dataType, analysisID))

	c.JSON(http.StatusOK, data)
}