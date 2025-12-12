package http

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadHandler struct{}

func NewUploadHandler(app *fiber.App) *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" && ext != ".gif" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file type. Only images are allowed."})
	}

	filename := uuid.New().String() + ext

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to open file"})
	}
	defer src.Close()

	fileBytes := make([]byte, file.Size)
	if _, err := src.Read(fileBytes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read file"})
	}
	supabaseURL := os.Getenv("SUPABASE_PROJECT_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if supabaseKey == "" {
		supabaseKey = os.Getenv("SUPABASE_ANON_KEY")
	}

	bucketName := c.FormValue("bucket")
	if bucketName == "" {
		bucketName = "media"
	}

	allowedBuckets := map[string]bool{
		"media":  true,
		"secure": true,
	}
	if !allowedBuckets[bucketName] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid bucket specified"})
	}

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

	fmt.Printf("DEBUG: ProjectURL='%s', DBURL='%s', KeyPresent=%v\n", supabaseURL, os.Getenv("SUPABASE_URL"), supabaseKey != "")
	fmt.Printf("DEBUG: Derived URL: %s\n", supabaseURL)

	if supabaseURL == "" || supabaseKey == "" {
		fmt.Println("DEBUG: Config missing return 500")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server storage configuration missing"})
	}

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", supabaseURL, bucketName, filename)
	req, err := http.NewRequest("POST", uploadURL, bytes.NewBuffer(fileBytes))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create upload request"})
	}

	req.Header.Set("Authorization", "Bearer "+supabaseKey)
	req.Header.Set("Content-Type", file.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload to storage"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Supabase Upload Error: %s\n", string(body))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Storage provider rejected upload"})
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", supabaseURL, bucketName, filename)

	return c.JSON(fiber.Map{
		"url": publicURL,
	})
}
