package repository

type StorageRepository interface {
	UploadFile(bucketName string, filePath string, data []byte, contentType string) (publicURL string, err error)
	MoveFile(bucketName string, srcPath string, destPath string) error
}
