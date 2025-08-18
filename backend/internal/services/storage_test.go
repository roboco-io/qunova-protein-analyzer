package services

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorageService_NewStorageService(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// Act
	service := NewStorageService(tempDir)

	// Assert
	assert.NotNil(t, service)
	assert.Equal(t, tempDir, service.uploadDir)
	assert.NotNil(t, service.fileMap)

	// Verify directory was created
	_, err := os.Stat(tempDir)
	assert.NoError(t, err)
}

func TestStorageService_SaveFile_ValidFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	service := NewStorageService(tempDir)
	
	fileContent := "HEADER    TEST PROTEIN\nATOM      1  N   ALA A   1      10.000  10.000  10.000\nEND"
	fileReader := strings.NewReader(fileContent)
	
	fileID := "test-file-id"
	originalName := "test.pdb"

	// Act
	savedPath, err := service.SaveFile(fileReader, fileID, originalName)

	// Assert
	require.NoError(t, err)
	assert.Contains(t, savedPath, fileID)
	assert.Contains(t, savedPath, ".pdb")

	// Verify file was actually saved
	_, err = os.Stat(savedPath)
	assert.NoError(t, err)

	// Verify file content
	content, err := os.ReadFile(savedPath)
	require.NoError(t, err)
	assert.Equal(t, fileContent, string(content))

	// Verify file info was stored
	info, err := service.GetFileInfo(fileID)
	require.NoError(t, err)
	assert.Equal(t, fileID, info.ID)
	assert.Equal(t, originalName, info.OriginalName)
	assert.Equal(t, savedPath, info.StoredPath)
}

func TestStorageService_SaveFile_EmptyFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	service := NewStorageService(tempDir)
	
	fileReader := strings.NewReader("")
	fileID := "empty-file"
	originalName := "empty.pdb"

	// Act
	savedPath, err := service.SaveFile(fileReader, fileID, originalName)

	// Assert
	require.NoError(t, err)
	
	// Verify empty file was saved
	info, err := os.Stat(savedPath)
	require.NoError(t, err)
	assert.Equal(t, int64(0), info.Size())
}

func TestStorageService_GetFilePath_ExistingFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	service := NewStorageService(tempDir)
	
	fileReader := strings.NewReader("test content")
	fileID := "existing-file"
	originalName := "test.pdb"
	
	savedPath, err := service.SaveFile(fileReader, fileID, originalName)
	require.NoError(t, err)

	// Act
	retrievedPath, err := service.GetFilePath(fileID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, savedPath, retrievedPath)
}

func TestStorageService_GetFilePath_NonExistentFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	service := NewStorageService(tempDir)

	// Act
	_, err := service.GetFilePath("non-existent-file")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

func TestStorageService_DeleteFile_ExistingFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	service := NewStorageService(tempDir)
	
	fileReader := strings.NewReader("test content")
	fileID := "file-to-delete"
	originalName := "test.pdb"
	
	savedPath, err := service.SaveFile(fileReader, fileID, originalName)
	require.NoError(t, err)

	// Verify file exists before deletion
	_, err = os.Stat(savedPath)
	require.NoError(t, err)

	// Act
	err = service.DeleteFile(fileID)

	// Assert
	require.NoError(t, err)

	// Verify file was deleted from filesystem
	_, err = os.Stat(savedPath)
	assert.True(t, os.IsNotExist(err))

	// Verify file info was removed from map
	_, err = service.GetFileInfo(fileID)
	assert.Error(t, err)
}

func TestStorageService_DeleteFile_NonExistentFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	service := NewStorageService(tempDir)

	// Act
	err := service.DeleteFile("non-existent-file")

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

func TestStorageService_GetFileUploadTime(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	service := NewStorageService(tempDir)
	
	before := time.Now()
	
	fileReader := strings.NewReader("test content")
	fileID := "timed-file"
	originalName := "test.pdb"
	
	_, err := service.SaveFile(fileReader, fileID, originalName)
	require.NoError(t, err)
	
	after := time.Now()

	// Act
	uploadTime := service.GetFileUploadTime(fileID)

	// Assert
	assert.True(t, uploadTime.After(before) || uploadTime.Equal(before))
	assert.True(t, uploadTime.Before(after) || uploadTime.Equal(after))
}

func TestStorageService_CleanupOldFiles(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	service := NewStorageService(tempDir)
	
	// Save some files
	fileReader1 := strings.NewReader("content1")
	fileReader2 := strings.NewReader("content2")
	
	_, err := service.SaveFile(fileReader1, "old-file", "old.pdb")
	require.NoError(t, err)
	
	_, err = service.SaveFile(fileReader2, "new-file", "new.pdb")
	require.NoError(t, err)
	
	// Manually set old upload time for testing
	if info, exists := service.fileMap["old-file"]; exists {
		info.UploadedAt = time.Now().Add(-2 * time.Hour) // 2 hours ago
		service.fileMap["old-file"] = info
	}

	// Act
	err = service.CleanupOldFiles(1 * time.Hour) // Clean files older than 1 hour

	// Assert
	require.NoError(t, err)
	
	// Old file should be deleted
	_, err = service.GetFilePath("old-file")
	assert.Error(t, err)
	
	// New file should still exist
	_, err = service.GetFilePath("new-file")
	assert.NoError(t, err)
}

func TestStorageService_FileExtensionHandling(t *testing.T) {
	tests := []struct {
		name         string
		originalName string
		expectedExt  string
	}{
		{"PDB file", "protein.pdb", ".pdb"},
		{"CIF file", "structure.cif", ".cif"},
		{"mmCIF file", "complex.mmcif", ".mmcif"},
		{"No extension", "noext", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tempDir := t.TempDir()
			service := NewStorageService(tempDir)
			fileReader := strings.NewReader("test content")

			// Act
			savedPath, err := service.SaveFile(fileReader, "test-id", tt.originalName)

			// Assert
			require.NoError(t, err)
			actualExt := filepath.Ext(savedPath)
			assert.Equal(t, tt.expectedExt, actualExt)
		})
	}
}