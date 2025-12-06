package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileService struct {
	UploadDir string
	MaxSize   int64 // in bytes
}

func NewFileService(uploadDir string, maxSizeMB int64) *FileService {
	// Create upload directory if not exists
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic("Failed to create upload directory: " + err.Error())
	}

	return &FileService{
		UploadDir: uploadDir,
		MaxSize:   maxSizeMB * 1024 * 1024, // Convert MB to bytes
	}
}

type UploadedFile struct {
	FileName     string `json:"fileName"`
	FileURL      string `json:"fileUrl"`
	FileType     string `json:"fileType"`
	FileSize     int64  `json:"fileSize"`
	OriginalName string `json:"originalName"`
}

// UploadFile - Upload single file
func (s *FileService) UploadFile(file *multipart.FileHeader) (*UploadedFile, error) {
	// Validate file size
	if file.Size > s.MaxSize {
		return nil, fmt.Errorf("file size exceeds maximum allowed size of %d MB", s.MaxSize/(1024*1024))
	}

	// Validate file type
	if err := s.validateFileType(file.Filename); err != nil {
		return nil, err
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	// Create full path
	filePath := filepath.Join(s.UploadDir, uniqueFilename)

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("failed to open uploaded file")
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, errors.New("failed to create destination file")
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, src); err != nil {
		return nil, errors.New("failed to save file")
	}

	// Return file info
	return &UploadedFile{
		FileName:     uniqueFilename,
		FileURL:      "/uploads/" + uniqueFilename,
		FileType:     s.getFileType(ext),
		FileSize:     file.Size,
		OriginalName: file.Filename,
	}, nil
}

// UploadMultipleFiles - Upload multiple files
func (s *FileService) UploadMultipleFiles(files []*multipart.FileHeader) ([]*UploadedFile, error) {
	uploadedFiles := make([]*UploadedFile, 0, len(files))

	for _, file := range files {
		uploaded, err := s.UploadFile(file)
		if err != nil {
			// Rollback: delete already uploaded files
			s.DeleteFiles(uploadedFiles)
			return nil, fmt.Errorf("failed to upload %s: %s", file.Filename, err.Error())
		}
		uploadedFiles = append(uploadedFiles, uploaded)
	}

	return uploadedFiles, nil
}

// DeleteFile - Delete single file
func (s *FileService) DeleteFile(filename string) error {
	filePath := filepath.Join(s.UploadDir, filename)
	return os.Remove(filePath)
}

// DeleteFiles - Delete multiple files
func (s *FileService) DeleteFiles(files []*UploadedFile) {
	for _, file := range files {
		_ = s.DeleteFile(file.FileName)
	}
}

// validateFileType - Validate allowed file types
func (s *FileService) validateFileType(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))

	allowedTypes := map[string]bool{
		".pdf":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
		".zip":  true,
		".rar":  true,
	}

	if !allowedTypes[ext] {
		return fmt.Errorf("file type %s is not allowed", ext)
	}

	return nil
}

// getFileType - Get file type category
func (s *FileService) getFileType(ext string) string {
	ext = strings.ToLower(ext)

	imageTypes := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	documentTypes := map[string]bool{
		".pdf":  true,
		".doc":  true,
		".docx": true,
		".xls":  true,
		".xlsx": true,
	}

	archiveTypes := map[string]bool{
		".zip": true,
		".rar": true,
	}

	if imageTypes[ext] {
		return "image"
	} else if documentTypes[ext] {
		return "document"
	} else if archiveTypes[ext] {
		return "archive"
	}

	return "other"
}
