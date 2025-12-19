package usecase

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/pur108/talestoon-be/internal/domain/repository"
)

type UploadUsecase interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}

type uploadUsecase struct {
	storageRepo repository.StorageRepository
}

func NewUploadUsecase(storageRepo repository.StorageRepository) UploadUsecase {
	return &uploadUsecase{storageRepo}
}
func (u *uploadUsecase) UploadFile(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", errors.New("no file uploaded")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".gif" {
		return "", errors.New("invalid file type. Only images are allowed")
	}

	// Security: Server enforces the bucket name. Client cannot specify it.
	bucketName := "media"

	return u.storageRepo.UploadFile(file, bucketName)
}
