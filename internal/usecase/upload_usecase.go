package usecase

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/pur108/webteen-be/internal/domain/repository"
)

type UploadResponse struct {
	URL string `json:"url"`
}

type UploadUsecase interface {
	UploadFile(userID string, filename string, data []byte, contentType string) (*UploadResponse, error)
}

type uploadUsecase struct {
	storageRepo repository.StorageRepository
}

func NewUploadUsecase(storageRepo repository.StorageRepository) UploadUsecase {
	return &uploadUsecase{storageRepo}
}
func (u *uploadUsecase) UploadFile(userID string, filename string, data []byte, contentType string) (*UploadResponse, error) {
	if filename == "" {
		return nil, errors.New("filename required")
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".gif" {
		return nil, errors.New("invalid file type. Only images are allowed")
	}

	// Validate file size (e.g. 5MB) - handled in handler or here?
	// For now assume handler passes valid size or we check len(data)
	if len(data) > 5*1024*1024 {
		return nil, errors.New("file too large. Max 5MB")
	}

	bucketName := "media"
	// Draft path: drafts/{userID}/{uuid}{ext}
	newFilename := fmt.Sprintf("drafts/%s/%s%s", userID, uuid.New().String(), ext)

	publicURL, err := u.storageRepo.UploadFile(bucketName, newFilename, data, contentType)
	if err != nil {
		return nil, err
	}

	return &UploadResponse{
		URL: publicURL,
	}, nil
}
