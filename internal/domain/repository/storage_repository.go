package repository

import "mime/multipart"

type StorageRepository interface {
	UploadFile(file *multipart.FileHeader, bucketName string) (string, error)
}
