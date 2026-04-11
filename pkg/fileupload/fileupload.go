package fileupload

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// FileUploader defines the contract for uploading and deleting files.
type FileUploader interface {
	Upload(file *multipart.FileHeader, subDir string, allowedExts []string) (string, error)
	Delete(fileURL string) error
}

type localFileUploader struct {
	baseDir    string // e.g., "./public/uploads"
	publicPath string // e.g., "/public/uploads"
}

// NewLocalFileUploader creates a new instance of LocalFileUploader.
func NewLocalFileUploader(baseDir, publicPath string) FileUploader {
	return &localFileUploader{
		baseDir:    baseDir,
		publicPath: publicPath,
	}
}

func (u *localFileUploader) Upload(file *multipart.FileHeader, subDir string, allowedExts []string) (string, error) {
	// 1. Validate extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if len(allowedExts) > 0 {
		isValid := false
		for _, e := range allowedExts {
			if strings.EqualFold(e, ext) {
				isValid = true
				break
			}
		}
		if !isValid {
			return "", fmt.Errorf("file extension %s is not allowed", ext)
		}
	}

	// 2. Prepare directory
	uploadDir := filepath.Join(u.baseDir, subDir)
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create upload directory: %w", err)
		}
	}

	// 3. Generate unique filename
	fileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, fileName)

	// 4. Save file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	// 5. Return normalized public URL
	// Result: /public/uploads/materials/uuid.jpg
	publicURL := "/" + filepath.ToSlash(filepath.Join(u.publicPath, subDir, fileName))
	return publicURL, nil
}

func (u *localFileUploader) Delete(fileURL string) error {
	// Remove leading slash for local path conversion
	relPath := strings.TrimPrefix(fileURL, "/")
	
	// Convert public path to local disk path
	// If publicURL is /public/uploads/xxx.jpg and publicPath is public/uploads
	// We need to map it back to baseDir
	
	// Implementation note: This assumes fileURL starts with u.publicPath
	// A more robust implementation would check this.
	
	diskPath := filepath.Join(u.baseDir, strings.TrimPrefix(relPath, strings.TrimPrefix(u.publicPath, "/")))
	
	if err := os.Remove(diskPath); err != nil {
		if os.IsNotExist(err) {
			return nil // Already gone
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
