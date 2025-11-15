package service

import (
	"github.com/gofiber/fiber/v2"
	"go_clean/app/models/mongodb"
	// "go_clean/app/repository/mongodb"
	"go_clean/utils"
)

type AuthMongoRepo interface {
    FindByUsernameOrEmail(identifier string) (*models.LoginMongo, error)
}

type AuthMongoService struct {
    Repo AuthMongoRepo
}


// Login godoc
// @Summary Login user
// @Description Mengembalikan JWT token
// @Tags Auth
// @Version 2.0
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Data"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /login-mongo [post]
func (s *AuthMongoService) Login(c *fiber.Ctx) error {
    var req models.LoginRequest

    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
    }

    user, err := s.Repo.FindByUsernameOrEmail(req.Username)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "username/password salah"})
    }

    // cek plain password vs hash
    if !utils.CheckPassword(req.Password, user.PasswordHash) {
        return c.Status(401).JSON(fiber.Map{"error": "username/password salah"})
    }

    token, err := utils.GenerateTokenMongo(*user)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "gagal membuat token"})
    }

    return c.JSON(fiber.Map{
        "user":  user,
        "token": token,
    })
}

