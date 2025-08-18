package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type StorageService struct {
	uploadDir string
	fileMap   map[string]FileInfo
}

type FileInfo struct {
	ID         string
	OriginalName string
	StoredPath  string
	UploadedAt  time.Time
}

func NewStorageService(uploadDir string) *StorageService {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create upload directory: %v", err))
	}

	return &StorageService{
		uploadDir: uploadDir,
		fileMap:   make(map[string]FileInfo),
	}
}

// SaveFile saves the uploaded file to disk
func (s *StorageService) SaveFile(file multipart.File, fileID, originalName string) (string, error) {
	// Create file path
	ext := filepath.Ext(originalName)
	fileName := fmt.Sprintf("%s%s", fileID, ext)
	filePath := filepath.Join(s.uploadDir, fileName)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Store file info
	s.fileMap[fileID] = FileInfo{
		ID:           fileID,
		OriginalName: originalName,
		StoredPath:   filePath,
		UploadedAt:   time.Now(),
	}

	return filePath, nil
}

// GetFilePath retrieves the file path for a given file ID
func (s *StorageService) GetFilePath(fileID string) (string, error) {
	info, exists := s.fileMap[fileID]
	if !exists {
		return "", fmt.Errorf("file not found: %s", fileID)
	}
	return info.StoredPath, nil
}

// GetFileInfo retrieves file information
func (s *StorageService) GetFileInfo(fileID string) (*FileInfo, error) {
	info, exists := s.fileMap[fileID]
	if !exists {
		return nil, fmt.Errorf("file not found: %s", fileID)
	}
	return &info, nil
}

// GetFileUploadTime retrieves the upload time for a file
func (s *StorageService) GetFileUploadTime(fileID string) time.Time {
	if info, exists := s.fileMap[fileID]; exists {
		return info.UploadedAt
	}
	return time.Now()
}

// DeleteFile removes a file from storage
func (s *StorageService) DeleteFile(fileID string) error {
	info, exists := s.fileMap[fileID]
	if !exists {
		return fmt.Errorf("file not found: %s", fileID)
	}

	// Remove file from disk
	if err := os.Remove(info.StoredPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	// Remove from map
	delete(s.fileMap, fileID)
	return nil
}

// CleanupOldFiles removes files older than the specified duration
func (s *StorageService) CleanupOldFiles(maxAge time.Duration) error {
	now := time.Now()
	var toDelete []string

	for fileID, info := range s.fileMap {
		if now.Sub(info.UploadedAt) > maxAge {
			toDelete = append(toDelete, fileID)
		}
	}

	for _, fileID := range toDelete {
		if err := s.DeleteFile(fileID); err != nil {
			return err
		}
	}

	return nil
}