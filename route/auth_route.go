package route

import (
	"go_clean/app/repository"
	"go_clean/app/service"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupAuthMongoRoutes(app *fiber.App, mongoDB *mongo.Database) {
	// 🔧 inisialisasi repository & service
	userRepo := repository.NewUserMongoRepository(mongoDB)
	authService := &service.AuthMongoService{
		Repo: userRepo,
	}

	// API Group tanpa middleware (login = public)
	api := app.Group("/api")

	// POST /api/login-mongo → login pakai MongoDB
	api.Post("/login-mongo", authService.Login)
}
