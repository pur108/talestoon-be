package repository

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/pur108/talestoon-be/internal/domain/repository"
)

type supabaseRepository struct{}

func NewSupabaseRepository() repository.StorageRepository {
	return &supabaseRepository{}
}

func (r *supabaseRepository) UploadFile(file *multipart.FileHeader, bucketName string) (string, error) {
	supabaseURL := os.Getenv("SUPABASE_PROJECT_URL")
	supabaseKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" {
		dbURL := os.Getenv("SUPABASE_URL")
		if parts := strings.Split(dbURL, "@"); len(parts) > 1 {
			if hostParts := strings.Split(parts[1], "."); len(hostParts) > 1 {
				if strings.Contains(dbURL, "supabase.co") {
					domainParts := strings.Split(parts[1], ":")
					if len(domainParts) > 0 {
						host := domainParts[0]
						ref := strings.TrimPrefix(host, "db.")
						supabaseURL = fmt.Sprintf("https://%s", ref)
					}
				}
			}
		}
	}

	var missingVars []string
	if supabaseURL == "" {
		missingVars = append(missingVars, "SUPABASE_PROJECT_URL")
	}
	if supabaseKey == "" {
		missingVars = append(missingVars, "SUPABASE_ANON_KEY")
	}

	if len(missingVars) > 0 {
		return "", fmt.Errorf("server storage configuration missing: %s", strings.Join(missingVars, ", "))
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Determine file extension
	ext := strings.ToLower(file.Filename[strings.LastIndex(file.Filename, "."):])
	filename := uuid.New().String() + ext

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucketName, filename)
	req, err := http.NewRequest("POST", uploadURL, bytes.NewBuffer(fileBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", file.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload to storage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("storage provider rejected upload: %s", string(body))
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucketName, filename)
	return publicURL, nil
}
