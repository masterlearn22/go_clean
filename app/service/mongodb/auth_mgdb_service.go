package service

import (
	"github.com/gofiber/fiber/v2"
	"go_clean/app/models/mongodb"
	"go_clean/app/repository/mongodb"
	"go_clean/utils"
)

type AuthMongoService struct {
	Repo *repository.UserMongoRepository
}

func (s *AuthMongoService) Login(c *fiber.Ctx) error {
	var req models.LoginMongo

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid payload"})
	}

	// ambil user
	user, err := s.Repo.FindByUsernameOrEmail(req.Username)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "username/password salah"})
	}

	// check password hash
	if !utils.CheckPassword(req.PasswordHash, user.PasswordHash) {
		return c.Status(401).JSON(fiber.Map{"error": "username/password salah"})
	}

	// generate jwt
	token, err := utils.GenerateTokenMongo(*user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal membuat token"})
	}

	return c.JSON(fiber.Map{
		"user":  user,
		"token": token,
	})
}
