package service

import (
    "strings"

	"github.com/gofiber/fiber/v2"
	"go_clean/app/models/postgresql"
	"go_clean/app/repository/postgresql"
	"go_clean/utils"
)

type AuthService struct {
    Repo *repository.AuthRepository
}

// Login godoc
// @Summary Login user
// @Description Mengembalikan JWT token
// @Tags Auth
// @Version 1.0
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Data"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /login [post]
func (s *AuthService) LoginUser(c *fiber.Ctx) error {
	var req models.LoginRequest

	// parse & validasi payload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "payload tidak valid"})
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	if req.Username == "" || req.Password == "" {
		return c.Status(400).JSON(fiber.Map{"error": "username & password wajib"})
	}

	// ambil data user dari repository
	u, hash, err := s.Repo.GetByUsernameOrEmail(req.Username)
	if err != nil {	
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "username/password salah",
		})
	}

	// cek password hash
	if !utils.CheckPassword(req.Password, hash) {
		return c.Status(401).JSON(fiber.Map{"error": "username/password salah"})
	}

	// generate token JWT
	token, err := utils.GenerateToken(*u)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal generate token"})
	}

	// return response model
	return c.JSON(models.LoginResponse{
		User:  *u,
		Token: token,
	})
}


