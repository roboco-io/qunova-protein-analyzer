package mocks

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/qunova/protein-analyzer/internal/models"
)

// MockStorageService provides a mock implementation for testing
type MockStorageService struct {
	files map[string]mockFileInfo
}

type mockFileInfo struct {
	ID           string
	OriginalName string
	StoredPath   string
	UploadedAt   time.Time
}

func NewMockStorageService() *MockStorageService {
	return &MockStorageService{
		files: make(map[string]mockFileInfo),
	}
}

func (m *MockStorageService) SaveFile(file multipart.File, fileID, originalName string) (string, error) {
	if fileID == "error-file" {
		return "", fmt.Errorf("mock save error")
	}
	
	filePath := fmt.Sprintf("/mock/path/%s%s", fileID, filepath.Ext(originalName))
	m.files[fileID] = mockFileInfo{
		ID:           fileID,
		OriginalName: originalName,
		StoredPath:   filePath,
		UploadedAt:   time.Now(),
	}
	
	return filePath, nil
}

func (m *MockStorageService) GetFilePath(fileID string) (string, error) {
	if fileID == "not-found" {
		return "", fmt.Errorf("file not found")
	}
	
	if info, exists := m.files[fileID]; exists {
		return info.StoredPath, nil
	}
	
	// Return default path for testing
	return fmt.Sprintf("/mock/path/%s.pdb", fileID), nil
}

func (m *MockStorageService) GetFileUploadTime(fileID string) time.Time {
	if info, exists := m.files[fileID]; exists {
		return info.UploadedAt
	}
	return time.Now()
}

func (m *MockStorageService) DeleteFile(fileID string) error {
	if fileID == "error-delete" {
		return fmt.Errorf("mock delete error")
	}
	
	delete(m.files, fileID)
	return nil
}

// MockParserService provides a mock implementation for testing
type MockParserService struct {
	ShouldError bool
	ErrorMsg    string
}

func NewMockParserService() *MockParserService {
	return &MockParserService{}
}

func (m *MockParserService) ParseStructureFile(filePath string) (*models.ProteinStructure, error) {
	if m.ShouldError {
		return nil, fmt.Errorf(m.ErrorMsg)
	}
	
	if filepath.Base(filePath) == "invalid.pdb" {
		return nil, fmt.Errorf("invalid file format")
	}
	
	return m.MockParseStructureFile(filePath), nil
}

func (m *MockParserService) MockParseStructureFile(filePath string) *models.ProteinStructure {
	return &models.ProteinStructure{
		ID:         "mock-structure-id",
		FileName:   filepath.Base(filePath),
		TotalAtoms: 100,
		Chains: []models.Chain{
			{ID: "A", Residues: 50, Atoms: 75, Type: "Protein"},
			{ID: "B", Residues: 25, Atoms: 25, Type: "DNA"},
		},
		Domains: []models.Domain{
			{Name: "Test Domain", Range: "1-50", Chains: []string{"A"}, Atoms: 75},
		},
		Ligands: []models.Ligand{
			{Name: "ATP", Atoms: 31, BindingSite: "Chain A:45-52"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// MockAnalysisService provides a mock implementation for testing
type MockAnalysisService struct {
	results map[string]*models.AnalysisResult
}

func NewMockAnalysisService() *MockAnalysisService {
	return &MockAnalysisService{
		results: make(map[string]*models.AnalysisResult),
	}
}

func (m *MockAnalysisService) ProcessStructure(analysisID, filePath string, options models.AnalysisOptions) {
	// Simulate immediate completion for testing
	now := time.Now()
	mockParser := NewMockParserService()
	structure := mockParser.MockParseStructureFile(filePath)
	
	result := &models.AnalysisResult{
		ID:          analysisID,
		FileID:      filePath,
		Status:      "completed",
		Progress:    100,
		Structure:   structure,
		StartedAt:   now,
		CompletedAt: &now,
	}
	
	m.results[analysisID] = result
}

func (m *MockAnalysisService) GetAnalysisResult(analysisID string) (*models.AnalysisResult, error) {
	if analysisID == "not-found" {
		return nil, fmt.Errorf("analysis not found")
	}
	
	if result, exists := m.results[analysisID]; exists {
		return result, nil
	}
	
	// Return a default result for testing
	now := time.Now()
	return &models.AnalysisResult{
		ID:          analysisID,
		Status:      "completed",
		Progress:    100,
		StartedAt:   now,
		CompletedAt: &now,
	}, nil
}