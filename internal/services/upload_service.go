package services

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/shapnesht/ecommerce-platform-go/internal/interfaces"
)

type UploadService struct {
	provider interfaces.UploadProvider
}

func (s *UploadService) UploadFile(productId uint, file *multipart.FileHeader) (string, error) {
	fileName := file.Filename
	ext := strings.ToLower(filepath.Ext(fileName))

	if !isValidImageExt(ext) {
		return "", fmt.Errorf("invalid file type :%s", ext)
	}

	path := fmt.Sprintf("/products/%d/%s", productId, fileName)

	return s.provider.UploadFile(file, path)
}

func (s *UploadService) DeleteFile(filePath string) error {
	return s.provider.DeleteFile(filePath)
}

func isValidImageExt(ext string) bool {
	validExts := []string{".jpg", ".jpeg", ".png", ".webp", ".gif"}

	for _, validExt := range validExts {
		if validExt == ext {
			return true
		}
	}
	return false
}
