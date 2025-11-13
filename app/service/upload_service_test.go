package service

import (
	"bytes"
	"fmt"
	"go_clean/app/repository"
	"mime/multipart"
	"net/textproto"
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/gofiber/fiber/v2"
)

func createTestFileRequest(t *testing.T, fieldName, fileName, contentType string, size int64) (*http.Request, *bytes.Buffer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Header multipart didefinisikan manual
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, fileName))
	h.Set("Content-Type", contentType)

	// Buat part dengan header yang benar
	part, err := writer.CreatePart(h)
	if err != nil {
		t.Fatal(err)
	}

	// Isi data dummy
	data := bytes.Repeat([]byte("A"), int(size))
	if _, err := part.Write(data); err != nil {
		t.Fatal(err)
	}

	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, body
}


func TestUploadFile_Success(t *testing.T) {
	app := fiber.New()
	mockRepo := repository.NewMockFileRepository()

	uploadPath := "./test_uploads"
	defer os.RemoveAll(uploadPath)

	service := NewFileService(mockRepo, uploadPath)

	app.Post("/upload", service.UploadFile)

	req, _ := createTestFileRequest(t, "file", "test.jpg", "image/jpeg", 1024)

	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusCreated {
		t.Errorf("expected 201, got %d", resp.StatusCode)
	}
}

func TestUploadFile_NoFile(t *testing.T) {
	app := fiber.New()
	mockRepo := repository.NewMockFileRepository()
	service := NewFileService(mockRepo, "./test_uploads")

	app.Post("/upload", service.UploadFile)

	req := httptest.NewRequest("POST", "/upload", nil)

	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUploadFile_InvalidType(t *testing.T) {
	app := fiber.New()
	mockRepo := repository.NewMockFileRepository()
	service := NewFileService(mockRepo, "./test_uploads")

	app.Post("/upload", service.UploadFile)

	req, _ := createTestFileRequest(t, "file", "test.exe", "application/octet-stream", 1024)

	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUploadFile_TooLarge(t *testing.T) {
	app := fiber.New(fiber.Config{
		BodyLimit: 20 * 1024 * 1024, // perbolehkan request besar
	})

	mockRepo := repository.NewMockFileRepository()
	service := NewFileService(mockRepo, "./test_uploads")

	// Middleware: override ukuran file SETELAH Fiber parsing request
	app.Use(func(c *fiber.Ctx) error {
		form, err := c.MultipartForm()
		if err == nil {
			for _, fhs := range form.File {
				for _, fh := range fhs {
					fh.Size = 11 * 1024 * 1024 // >10MB
				}
			}
		}
		return c.Next()
	})

	app.Post("/upload", service.UploadFile)

	// request normal 1KB
	req, _ := createTestFileRequest(t, "file", "big.jpg", "image/jpeg", 1024)

	resp, _ := app.Test(req)
	if resp.StatusCode != fiber.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestUploadFile_DBError(t *testing.T) {
	app := fiber.New()
	mockRepo := repository.NewMockFileRepository()
	mockRepo.ForceError = true

	service := NewFileService(mockRepo, "./test_uploads")

	app.Post("/upload", service.UploadFile)

	req, _ := createTestFileRequest(t, "file", "test.jpg", "image/jpeg", 1024)

	resp, _ := app.Test(req)

	if resp.StatusCode != fiber.StatusInternalServerError {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}
