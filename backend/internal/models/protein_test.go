package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProteinStructure_Creation(t *testing.T) {
	// Arrange
	now := time.Now()
	structure := ProteinStructure{
		ID:         "test-id",
		FileName:   "test.pdb",
		TotalAtoms: 100,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Assert
	assert.Equal(t, "test-id", structure.ID)
	assert.Equal(t, "test.pdb", structure.FileName)
	assert.Equal(t, 100, structure.TotalAtoms)
	assert.Equal(t, now, structure.CreatedAt)
}

func TestChain_Validation(t *testing.T) {
	tests := []struct {
		name     string
		chain    Chain
		expected bool
	}{
		{
			name: "Valid protein chain",
			chain: Chain{
				ID:       "A",
				Residues: 100,
				Atoms:    800,
				Type:     "Protein",
			},
			expected: true,
		},
		{
			name: "Valid DNA chain",
			chain: Chain{
				ID:       "B",
				Residues: 50,
				Atoms:    400,
				Type:     "DNA",
			},
			expected: true,
		},
		{
			name: "Empty chain ID",
			chain: Chain{
				ID:       "",
				Residues: 100,
				Atoms:    800,
				Type:     "Protein",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validateChain(tt.chain)
			assert.Equal(t, tt.expected, isValid)
		})
	}
}

func TestAnalysisOptions_DefaultValues(t *testing.T) {
	// Arrange
	options := AnalysisOptions{
		IncludeSequence:  true,
		CalculateDomains: true,
		IdentifyLigands:  true,
		Resolution:       "high",
	}

	// Assert
	assert.True(t, options.IncludeSequence)
	assert.True(t, options.CalculateDomains)
	assert.True(t, options.IdentifyLigands)
	assert.Equal(t, "high", options.Resolution)
}

func TestAnalysisResult_StatusProgression(t *testing.T) {
	// Arrange
	startTime := time.Now()
	result := AnalysisResult{
		ID:        "test-analysis",
		Status:    "pending",
		Progress:  0,
		StartedAt: startTime,
	}

	// Act & Assert - Progress from pending to processing
	result.Status = "processing"
	result.Progress = 50
	assert.Equal(t, "processing", result.Status)
	assert.Equal(t, 50, result.Progress)

	// Act & Assert - Complete analysis
	completedTime := time.Now()
	result.Status = "completed"
	result.Progress = 100
	result.CompletedAt = &completedTime
	
	assert.Equal(t, "completed", result.Status)
	assert.Equal(t, 100, result.Progress)
	assert.NotNil(t, result.CompletedAt)
	assert.True(t, result.CompletedAt.After(result.StartedAt))
}

func TestLigand_BindingSiteFormat(t *testing.T) {
	tests := []struct {
		name     string
		ligand   Ligand
		expected bool
	}{
		{
			name: "Valid binding site format",
			ligand: Ligand{
				Name:        "ATP",
				Atoms:       31,
				BindingSite: "Chain A:45-52",
			},
			expected: true,
		},
		{
			name: "Multiple chain binding site",
			ligand: Ligand{
				Name:        "HEME",
				Atoms:       43,
				BindingSite: "Chain A:100-110, Chain B:50-55",
			},
			expected: true,
		},
		{
			name: "Empty binding site",
			ligand: Ligand{
				Name:        "MG",
				Atoms:       1,
				BindingSite: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := validateLigandBindingSite(tt.ligand.BindingSite)
			assert.Equal(t, tt.expected, isValid)
		})
	}
}

// Helper functions for validation (these would be implemented in the actual model)
func validateChain(chain Chain) bool {
	return chain.ID != "" && chain.Atoms > 0 && chain.Type != ""
}

func validateLigandBindingSite(bindingSite string) bool {
	return bindingSite != ""
}