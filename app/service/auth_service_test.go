package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"go_clean/app/models"
	"go_clean/app/repository"
	"go_clean/utils"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func resetMocks() {
	utils.MockCheckPassword = nil
	utils.MockGenerateToken = nil
}

func TestLogin_Success(t *testing.T) {
	resetMocks()
	app := fiber.New()

	mockRepo := repository.NewMockUserRepository()
	mockRepo.User = &models.User{ID: 1, Username: "admin", Email: "admin@mail.com", Role: "admin"}
	mockRepo.Hash = "$2a$10$abcdefghijklmnopqrstuv" // dummy hash

	utils.MockCheckPassword = func(pw, hash string) bool { return true }
	utils.MockGenerateToken = func(u models.User) (string, error) { return "token123", nil }

	service := &UserService{Repo: mockRepo}

	app.Post("/login", service.LoginUser)

	body := map[string]string{"username": "admin", "password": "123456"}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestLogin_InvalidPayload(t *testing.T) {
	resetMocks()
	app := fiber.New()

	mockRepo := repository.NewMockUserRepository()
	service := &UserService{Repo: mockRepo}

	app.Post("/login", service.LoginUser)

	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte("{invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestLogin_EmptyFields(t *testing.T) {
	resetMocks()
	app := fiber.New()

	mockRepo := repository.NewMockUserRepository()
	service := &UserService{Repo: mockRepo}

	app.Post("/login", service.LoginUser)

	body := `{"username": "", "password": ""}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 400 {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	resetMocks()
	app := fiber.New()

	mockRepo := repository.NewMockUserRepository()
	mockRepo.Err = errors.New("not found")

	service := &UserService{Repo: mockRepo}

	app.Post("/login", service.LoginUser)

	body := `{"username":"unknown","password":"123"}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	resetMocks()
	app := fiber.New()

	mockRepo := repository.NewMockUserRepository()
	mockRepo.User = &models.User{ID: 1, Username: "admin"}
	mockRepo.Hash = "correcthash"

	utils.MockCheckPassword = func(pw, hash string) bool { return false }

	service := &UserService{Repo: mockRepo}

	app.Post("/login", service.LoginUser)

	body := `{"username":"admin","password":"wrong"}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 401 {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestLogin_GenerateTokenError(t *testing.T) {
	resetMocks()
	app := fiber.New()

	mockRepo := repository.NewMockUserRepository()
	mockRepo.User = &models.User{ID: 1, Username: "admin"}
	mockRepo.Hash = "hash"

	utils.MockCheckPassword = func(pw, hash string) bool { return true }
	utils.MockGenerateToken = func(u models.User) (string, error) { return "", errors.New("token error") }

	service := &UserService{Repo: mockRepo}

	app.Post("/login", service.LoginUser)

	body := `{"username":"admin","password":"123"}`
	req := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	if resp.StatusCode != 500 {
		t.Errorf("expected 500, got %d", resp.StatusCode)
	}
}
