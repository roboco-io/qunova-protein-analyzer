package services

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/qunova/protein-analyzer/internal/models"
)

type ParserService struct {
	pythonScriptPath string
}

func NewParserService() *ParserService {
	return &ParserService{
		pythonScriptPath: "./scripts/parse_protein.py",
	}
}

// ParseStructureFile parses a protein structure file using Python script
func (s *ParserService) ParseStructureFile(filePath string) (*models.ProteinStructure, error) {
	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Execute Python parser script
	cmd := exec.Command("python3", s.pythonScriptPath, absPath)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	// Parse JSON output from Python script
	var structure models.ProteinStructure
	if err := json.Unmarshal(output, &structure); err != nil {
		return nil, fmt.Errorf("failed to parse parser output: %w", err)
	}

	return &structure, nil
}

// MockParseStructureFile provides mock data for testing without Python dependency
func (s *ParserService) MockParseStructureFile(filePath string) *models.ProteinStructure {
	// Return mock data for demonstration
	return &models.ProteinStructure{
		FileName:   filepath.Base(filePath),
		TotalAtoms: 2847,
		Chains: []models.Chain{
			{ID: "A", Residues: 142, Atoms: 1123, Type: "Protein"},
			{ID: "B", Residues: 89, Atoms: 712, Type: "Protein"},
			{ID: "C", Residues: 45, Atoms: 356, Type: "DNA"},
			{ID: "W", Residues: 127, Atoms: 381, Type: "Water"},
		},
		Domains: []models.Domain{
			{Name: "N-terminal domain", Range: "1-67", Chains: []string{"A"}, Atoms: 534},
			{Name: "C-terminal domain", Range: "68-142", Chains: []string{"A"}, Atoms: 589},
			{Name: "Binding domain", Range: "1-89", Chains: []string{"B"}, Atoms: 712},
		},
		Ligands: []models.Ligand{
			{Name: "ATP", Atoms: 31, BindingSite: "Chain A:45-52"},
			{Name: "MG", Atoms: 1, BindingSite: "Chain A:48"},
		},
	}
}