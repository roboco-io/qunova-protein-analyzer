package models

import (
	"time"
)

// ProteinStructure represents the complete protein structure data
type ProteinStructure struct {
	ID         string    `json:"id"`
	FileName   string    `json:"fileName"`
	TotalAtoms int       `json:"totalAtoms"`
	Chains     []Chain   `json:"chains"`
	Domains    []Domain  `json:"domains"`
	Ligands    []Ligand  `json:"ligands"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// Chain represents a protein chain
type Chain struct {
	ID       string `json:"id"`
	Residues int    `json:"residues"`
	Atoms    int    `json:"atoms"`
	Type     string `json:"type"` // Protein, DNA, RNA, Water
	Sequence string `json:"sequence,omitempty"`
}

// Domain represents a protein domain
type Domain struct {
	Name   string   `json:"name"`
	Range  string   `json:"range"`
	Chains []string `json:"chains"`
	Atoms  int      `json:"atoms"`
}

// Ligand represents a ligand molecule
type Ligand struct {
	Name        string `json:"name"`
	Atoms       int    `json:"atoms"`
	BindingSite string `json:"binding_site"`
}

// AnalysisRequest represents a request to analyze a protein structure
type AnalysisRequest struct {
	FileID  string            `json:"fileId"`
	Options AnalysisOptions   `json:"options,omitempty"`
}

// AnalysisOptions represents optional analysis parameters
type AnalysisOptions struct {
	IncludeSequence   bool `json:"includeSequence"`
	CalculateDomains  bool `json:"calculateDomains"`
	IdentifyLigands   bool `json:"identifyLigands"`
	Resolution        string `json:"resolution"` // low, medium, high
}

// AnalysisResult represents the result of protein structure analysis
type AnalysisResult struct {
	ID           string           `json:"id"`
	FileID       string           `json:"fileId"`
	Status       string           `json:"status"` // pending, processing, completed, failed
	Progress     int              `json:"progress"`
	Structure    *ProteinStructure `json:"structure,omitempty"`
	Error        string           `json:"error,omitempty"`
	StartedAt    time.Time        `json:"startedAt"`
	CompletedAt  *time.Time       `json:"completedAt,omitempty"`
}

// UploadResponse represents the response after file upload
type UploadResponse struct {
	FileID    string    `json:"fileId"`
	FileName  string    `json:"fileName"`
	FileSize  int64     `json:"fileSize"`
	FileType  string    `json:"fileType"`
	UploadedAt time.Time `json:"uploadedAt"`
}

// ExportRequest represents a request to export data
type ExportRequest struct {
	AnalysisID string `json:"analysisId"`
	DataType   string `json:"dataType"` // all, chains, domains, ligands
	Format     string `json:"format"`   // json, csv, pdb
}