package interfaces

import (
	"mime/multipart"
	"time"

	"github.com/qunova/protein-analyzer/internal/models"
)

// StorageInterface defines the contract for file storage operations
type StorageInterface interface {
	SaveFile(file multipart.File, fileID, originalName string) (string, error)
	GetFilePath(fileID string) (string, error)
	GetFileUploadTime(fileID string) time.Time
	DeleteFile(fileID string) error
}

// ParserInterface defines the contract for protein structure parsing
type ParserInterface interface {
	ParseStructureFile(filePath string) (*models.ProteinStructure, error)
	MockParseStructureFile(filePath string) *models.ProteinStructure
}

// AnalysisInterface defines the contract for protein analysis operations
type AnalysisInterface interface {
	ProcessStructure(analysisID, filePath string, options models.AnalysisOptions)
	GetAnalysisResult(analysisID string) (*models.AnalysisResult, error)
}