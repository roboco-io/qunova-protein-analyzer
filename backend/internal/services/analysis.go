package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/qunova/protein-analyzer/internal/models"
)

type AnalysisService struct {
	parserService *ParserService
	analysisMap   map[string]*models.AnalysisResult
	mu            sync.RWMutex
}

func NewAnalysisService(parserService *ParserService) *AnalysisService {
	return &AnalysisService{
		parserService: parserService,
		analysisMap:   make(map[string]*models.AnalysisResult),
	}
}

// ProcessStructure processes a protein structure file
func (s *AnalysisService) ProcessStructure(analysisID, filePath string, options models.AnalysisOptions) {
	// Create initial analysis result
	startTime := time.Now()
	result := &models.AnalysisResult{
		ID:        analysisID,
		FileID:    filePath,
		Status:    "processing",
		Progress:  0,
		StartedAt: startTime,
	}

	// Store initial result
	s.mu.Lock()
	s.analysisMap[analysisID] = result
	s.mu.Unlock()

	// Simulate progress updates
	go s.updateProgress(analysisID)

	// Parse structure file (using mock for now)
	structure := s.parserService.MockParseStructureFile(filePath)
	
	// In production, use actual parser:
	// structure, err := s.parserService.ParseStructureFile(filePath)
	// if err != nil {
	//     s.updateAnalysisError(analysisID, err)
	//     return
	// }

	// Add analysis-specific processing based on options
	if options.CalculateDomains {
		// Domain calculation logic would go here
		structure.Domains = s.calculateDomains(structure)
	}

	if options.IdentifyLigands {
		// Ligand identification logic would go here
		structure.Ligands = s.identifyLigands(structure)
	}

	// Set structure ID
	structure.ID = uuid.New().String()
	structure.CreatedAt = time.Now()
	structure.UpdatedAt = time.Now()

	// Update result with completed analysis
	completedTime := time.Now()
	s.mu.Lock()
	result.Status = "completed"
	result.Progress = 100
	result.Structure = structure
	result.CompletedAt = &completedTime
	s.analysisMap[analysisID] = result
	s.mu.Unlock()
}

// GetAnalysisResult retrieves an analysis result
func (s *AnalysisService) GetAnalysisResult(analysisID string) (*models.AnalysisResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, exists := s.analysisMap[analysisID]
	if !exists {
		return nil, fmt.Errorf("analysis not found: %s", analysisID)
	}

	return result, nil
}

// updateProgress simulates progress updates
func (s *AnalysisService) updateProgress(analysisID string) {
	progress := 0
	for progress < 90 {
		time.Sleep(500 * time.Millisecond)
		progress += 10

		s.mu.Lock()
		if result, exists := s.analysisMap[analysisID]; exists {
			result.Progress = progress
		}
		s.mu.Unlock()
	}
}

// updateAnalysisError updates analysis with error status
func (s *AnalysisService) updateAnalysisError(analysisID string, err error) {
	completedTime := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()

	if result, exists := s.analysisMap[analysisID]; exists {
		result.Status = "failed"
		result.Error = err.Error()
		result.CompletedAt = &completedTime
	}
}

// calculateDomains calculates protein domains (placeholder)
func (s *AnalysisService) calculateDomains(structure *models.ProteinStructure) []models.Domain {
	// This is a placeholder. In production, implement actual domain calculation
	if structure.Domains == nil || len(structure.Domains) == 0 {
		return []models.Domain{
			{
				Name:   "Default domain",
				Range:  "1-100",
				Chains: []string{"A"},
				Atoms:  1000,
			},
		}
	}
	return structure.Domains
}

// identifyLigands identifies ligands in the structure (placeholder)
func (s *AnalysisService) identifyLigands(structure *models.ProteinStructure) []models.Ligand {
	// This is a placeholder. In production, implement actual ligand identification
	if structure.Ligands == nil || len(structure.Ligands) == 0 {
		return []models.Ligand{
			{
				Name:        "HEM",
				Atoms:       43,
				BindingSite: "Chain A:100-110",
			},
		}
	}
	return structure.Ligands
}

// CleanupOldAnalyses removes analyses older than the specified duration
func (s *AnalysisService) CleanupOldAnalyses(maxAge time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	var toDelete []string

	for id, result := range s.analysisMap {
		if result.CompletedAt != nil && now.Sub(*result.CompletedAt) > maxAge {
			toDelete = append(toDelete, id)
		}
	}

	for _, id := range toDelete {
		delete(s.analysisMap, id)
	}
}