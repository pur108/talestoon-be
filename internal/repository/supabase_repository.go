package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/pur108/talestoon-be/internal/domain/repository"
)

type supabaseRepository struct{}

func NewSupabaseRepository() repository.StorageRepository {
	return &supabaseRepository{}
}

func (r *supabaseRepository) UploadFile(bucketName string, filePath string, data []byte, contentType string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_PROJECT_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL == "" {
		return "", fmt.Errorf("server storage configuration missing: SUPABASE_PROJECT_URL")
	}
	if supabaseKey == "" {
		return "", fmt.Errorf("server storage configuration missing: SUPABASE_SERVICE_ROLE_KEY")
	}

	// Upload endpoint: POST /storage/v1/object/{bucket}/{wildcard}
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucketName, filePath)

	req, err := http.NewRequest("POST", uploadURL, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("storage provider rejected upload: %s", string(body))
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucketName, filePath)
	return publicURL, nil
}

func (r *supabaseRepository) MoveFile(bucketName string, srcPath string, destPath string) error {
	supabaseURL := os.Getenv("SUPABASE_PROJECT_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL == "" {
		return fmt.Errorf("server storage configuration missing: SUPABASE_PROJECT_URL")
	}
	if supabaseKey == "" {
		return fmt.Errorf("server storage configuration missing: SUPABASE_SERVICE_ROLE_KEY")
	}

	moveURL := fmt.Sprintf("%s/storage/v1/object/move", supabaseURL)

	payload := map[string]string{
		"bucketId":       bucketName,
		"sourceKey":      srcPath,
		"destinationKey": destPath,
	}
	jsonBody, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", moveURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create move request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("storage provider rejected move: %s", string(body))
	}

	return nil
}
