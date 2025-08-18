package services

import (
	"testing"
	"time"

	"github.com/qunova/protein-analyzer/internal/mocks"
	"github.com/qunova/protein-analyzer/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnalysisService_NewAnalysisService(t *testing.T) {
	// Arrange
	mockParser := mocks.NewMockParserService()

	// Act
	service := NewAnalysisService(mockParser)

	// Assert
	assert.NotNil(t, service)
	assert.NotNil(t, service.parserService)
	assert.NotNil(t, service.analysisMap)
}

func TestAnalysisService_ProcessStructure_Success(t *testing.T) {
	// Arrange
	mockParser := mocks.NewMockParserService()
	service := NewAnalysisService(mockParser)
	
	analysisID := "test-analysis-id"
	filePath := "test.pdb"
	options := models.AnalysisOptions{
		IncludeSequence:  true,
		CalculateDomains: true,
		IdentifyLigands:  true,
		Resolution:       "high",
	}

	// Act
	service.ProcessStructure(analysisID, filePath, options)

	// Give some time for async processing
	time.Sleep(100 * time.Millisecond)

	// Assert
	result, err := service.GetAnalysisResult(analysisID)
	require.NoError(t, err)
	assert.Equal(t, analysisID, result.ID)
	assert.Equal(t, "completed", result.Status)
	assert.Equal(t, 100, result.Progress)
	assert.NotNil(t, result.Structure)
	assert.NotNil(t, result.CompletedAt)
}

func TestAnalysisService_ProcessStructure_WithOptions(t *testing.T) {
	// Arrange
	mockParser := mocks.NewMockParserService()
	service := NewAnalysisService(mockParser)
	
	analysisID := "test-with-options"
	filePath := "complex.pdb"
	options := models.AnalysisOptions{
		IncludeSequence:  false,
		CalculateDomains: true,
		IdentifyLigands:  false,
		Resolution:       "low",
	}

	// Act
	service.ProcessStructure(analysisID, filePath, options)
	time.Sleep(100 * time.Millisecond)

	// Assert
	result, err := service.GetAnalysisResult(analysisID)
	require.NoError(t, err)
	
	// Verify that domains were calculated (CalculateDomains = true)
	assert.NotEmpty(t, result.Structure.Domains)
	
	// In a real implementation, we would check that ligands were not identified
	// based on the IdentifyLigands = false option
}

func TestAnalysisService_GetAnalysisResult_ExistingAnalysis(t *testing.T) {
	// Arrange
	mockParser := mocks.NewMockParserService()
	service := NewAnalysisService(mockParser)
	
	analysisID := "existing-analysis"
	filePath := "test.pdb"
	options := models.AnalysisOptions{}

	// Create an analysis
	service.ProcessStructure(analysisID, filePath, options)
	time.Sleep(100 * time.Millisecond)

	// Act
	result, err := service.GetAnalysisResult(analysisID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, analysisID, result.ID)
	assert.NotEmpty(t, result.Status)
}

func TestAnalysisService_GetAnalysisResult_NonExistentAnalysis(t *testing.T) {
	// Arrange
	mockParser := mocks.NewMockParserService()
	service := NewAnalysisService(mockParser)

	// Act
	result, err := service.GetAnalysisResult("non-existent-id")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "analysis not found")
}

func TestAnalysisService_CalculateDomains(t *testing.T) {
	// Arrange
	mockParser := mocks.NewMockParserService()
	service := NewAnalysisService(mockParser)
	
	structure := &models.ProteinStructure{
		FileName:   "test.pdb",
		TotalAtoms: 1000,
		Chains: []models.Chain{
			{ID: "A", Residues: 100, Atoms: 800, Type: "Protein"},
		},
	}

	// Act
	domains := service.calculateDomains(structure)

	// Assert
	assert.NotEmpty(t, domains)
	// Should create at least one domain for the test
	assert.GreaterOrEqual(t, len(domains), 1)
	
	for _, domain := range domains {
		assert.NotEmpty(t, domain.Name)
		assert.NotEmpty(t, domain.Range)
		assert.NotEmpty(t, domain.Chains)
		assert.Greater(t, domain.Atoms, 0)
	}
}

func TestAnalysisService_IdentifyLigands(t *testing.T) {
	// Arrange
	mockParser := mocks.NewMockParserService()
	service := NewAnalysisService(mockParser)
	
	structure := &models.ProteinStructure{
		FileName:   "test.pdb",
		TotalAtoms: 1000,
		Ligands:    nil, // Start with no ligands
	}

	// Act
	ligands := service.identifyLigands(structure)

	// Assert
	assert.NotEmpty(t, ligands)
	
	for _, ligand := range ligands {
		assert.NotEmpty(t, ligand.Name)
		assert.Greater(t, ligand.Atoms, 0)
		assert.NotEmpty(t, ligand.BindingSite)
	}
}

func TestAnalysisService_CleanupOldAnalyses(t *testing.T) {
	// Arrange
	mockParser := mocks.NewMockParserService()
	service := NewAnalysisService(mockParser)
	
	// Create some analyses
	oldTime := time.Now().Add(-2 * time.Hour)
	newTime := time.Now()
	
	service.analysisMap["old-analysis"] = &models.AnalysisResult{
		ID:          "old-analysis",
		Status:      "completed",
		CompletedAt: &oldTime,
	}
	
	service.analysisMap["new-analysis"] = &models.AnalysisResult{
		ID:          "new-analysis", 
		Status:      "completed",
		CompletedAt: &newTime,
	}
	
	service.analysisMap["pending-analysis"] = &models.AnalysisResult{
		ID:          "pending-analysis",
		Status:      "processing",
		CompletedAt: nil, // Not completed yet
	}

	// Act
	service.CleanupOldAnalyses(1 * time.Hour) // Clean analyses older than 1 hour

	// Assert
	// Old completed analysis should be removed
	_, exists := service.analysisMap["old-analysis"]
	assert.False(t, exists)
	
	// New completed analysis should remain
	_, exists = service.analysisMap["new-analysis"]
	assert.True(t, exists)
	
	// Pending analysis should remain (not completed)
	_, exists = service.analysisMap["pending-analysis"]
	assert.True(t, exists)
}

func TestAnalysisService_ConcurrentAccess(t *testing.T) {
	// Arrange
	mockParser := mocks.NewMockParserService()
	service := NewAnalysisService(mockParser)
	
	analysisCount := 10
	done := make(chan bool, analysisCount)

	// Act - Start multiple analyses concurrently
	for i := 0; i < analysisCount; i++ {
		go func(id int) {
			analysisID := fmt.Sprintf("concurrent-analysis-%d", id)
			service.ProcessStructure(analysisID, "test.pdb", models.AnalysisOptions{})
			done <- true
		}(i)
	}

	// Wait for all to complete
	for i := 0; i < analysisCount; i++ {
		<-done
	}

	time.Sleep(200 * time.Millisecond) // Allow processing to complete

	// Assert - All analyses should be accessible
	for i := 0; i < analysisCount; i++ {
		analysisID := fmt.Sprintf("concurrent-analysis-%d", i)
		result, err := service.GetAnalysisResult(analysisID)
		assert.NoError(t, err)
		assert.Equal(t, analysisID, result.ID)
	}
}

func TestAnalysisService_ProgressTracking(t *testing.T) {
	// Arrange
	mockParser := mocks.NewMockParserService()
	service := NewAnalysisService(mockParser)
	
	analysisID := "progress-test"

	// Act
	service.ProcessStructure(analysisID, "test.pdb", models.AnalysisOptions{})

	// Check progress multiple times
	var progressValues []int
	for i := 0; i < 5; i++ {
		time.Sleep(50 * time.Millisecond)
		result, err := service.GetAnalysisResult(analysisID)
		if err == nil {
			progressValues = append(progressValues, result.Progress)
		}
	}

	// Assert
	assert.NotEmpty(t, progressValues)
	
	// Final progress should be 100
	finalResult, err := service.GetAnalysisResult(analysisID)
	require.NoError(t, err)
	assert.Equal(t, 100, finalResult.Progress)
	assert.Equal(t, "completed", finalResult.Status)
}