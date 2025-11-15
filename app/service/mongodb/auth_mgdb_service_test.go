package service

import (
    "testing"
    "go_clean/app/models/mongodb"
    "go_clean/app/repository/mongodb"
    "go_clean/utils"
    "net/http/httptest"
    "net/http"
    "encoding/json"
    "strings"
    "github.com/gofiber/fiber/v2"
)


func newTestRequest(app *fiber.App, route string, body interface{}) (*http.Response, error) {
    jsonBody, _ := json.Marshal(body)
    req := httptest.NewRequest("POST", route, strings.NewReader(string(jsonBody)))
    req.Header.Set("Content-Type", "application/json")
    return app.Test(req)
}

func TestLoginSuccess(t *testing.T) {
    mockRepo := repository.NewMockUserMongoRepository()
    svc := &AuthMongoService{Repo: mockRepo}

    hash, _ := utils.HashPassword("123456")

    mockRepo.InsertUser(&models.LoginMongo{
        Username:     "admin",
        Email:        "admin@mail.com",
        PasswordHash: hash,
        Role:         "admin",
    })

    // mock token generator
    utils.MockGenerateTokenMongo = func(u models.LoginMongo) (string, error) {
        return "mock-token", nil
    }
    defer func() { utils.MockGenerateTokenMongo = nil }()

    app := fiber.New()
    app.Post("/login-mongo", svc.Login)

    body := models.LoginRequest{Username: "admin", Password: "123456"}

    resp, err := newTestRequest(app, "/login-mongo", body)
    if err != nil {
        t.Errorf("request error: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected 200, got %v", resp.StatusCode)
    }
}

func TestLoginWrongPassword(t *testing.T) {
    mockRepo := repository.NewMockUserMongoRepository()
    svc := &AuthMongoService{Repo: mockRepo}

    hash, _ := utils.HashPassword("correctpass")

    mockRepo.InsertUser(&models.LoginMongo{
        Username:     "admin",
        PasswordHash: hash,
    })

    app := fiber.New()
    app.Post("/login-mongo", svc.Login)

    body := models.LoginRequest{Username: "admin", Password: "wrongpass"}

    resp, _ := newTestRequest(app, "/login-mongo", body)

    if resp.StatusCode != 401 {
        t.Errorf("expected 401 wrong password, got %v", resp.StatusCode)
    }
}

func TestLoginUserNotFound(t *testing.T) {
    mockRepo := repository.NewMockUserMongoRepository()
    svc := &AuthMongoService{Repo: mockRepo}

    app := fiber.New()
    app.Post("/login-mongo", svc.Login)

    body := models.LoginRequest{Username: "ghost", Password: "123"}

    resp, _ := newTestRequest(app, "/login-mongo", body)

    if resp.StatusCode != 401 {
        t.Errorf("expected 401 user not found, got %v", resp.StatusCode)
    }
}

func TestLoginTokenError(t *testing.T) {
    mockRepo := repository.NewMockUserMongoRepository()
    svc := &AuthMongoService{Repo: mockRepo}

    hash, _ := utils.HashPassword("123456")

    mockRepo.InsertUser(&models.LoginMongo{
        Username:     "admin",
        PasswordHash: hash,
    })

    // token generator bermasalah
    utils.MockGenerateTokenMongo = func(u models.LoginMongo) (string, error) {
        return "", fiber.ErrInternalServerError
    }
    defer func() { utils.MockGenerateTokenMongo = nil }()

    app := fiber.New()
    app.Post("/login-mongo", svc.Login)

    body := models.LoginRequest{Username: "admin", Password: "123456"}

    resp, _ := newTestRequest(app, "/login-mongo", body)

    if resp.StatusCode != 500 {
        t.Errorf("expected 500 token error, got %v", resp.StatusCode)
    }
}
