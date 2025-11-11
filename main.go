package main

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @tag.name Auth
// @tag.description Endpoint untuk autentikasi dan login user
// @tag.order 1

// @tag.name Alumni-Mongo
// @tag.description Endpoint untuk data alumni
// @tag.order 2

// @tag.name Pekerjaan-Mongo
// @tag.description Endpoint untuk data pekerjaan alumni (MongoDB)
// @tag.order 3

// @tag.name FileUpload
// @tag.description Upload, lihat dan hapus file
// @tag.order 4

// @tag.name Alumni-PostgresSQL
// @tag.description Upload, lihat dan hapus file
// @tag.order 5






import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go_clean/config"
	"go_clean/database"
	"go_clean/route"
	"go_clean/app/repository"
	"go_clean/app/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	 fiberSwagger "github.com/swaggo/fiber-swagger"

	docs "go_clean/docs"
	_ "go_clean/docs"
)

// @title User API
// @version 1.0
// @description API untuk mengelola data user dengan MongoDB menggunakan Clean Architecture
// @host localhost:3000
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization


func main() {
	// 1️ Load env
	config.LoadEnv()

	// 2️ Connect ke PostgreSQL
	database.ConnectDB()
	defer database.DB.Close()

	// 3️ Connect ke MongoDB
	database.ConnectMongoDB()

	// 4️ Setup Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error":   true,
				"message": err.Error(),
			})
		},
	})

	docs.SwaggerInfo.BasePath = "/api"
    app.Get("/swagger/*", fiberSwagger.WrapHandler)
	log.Println("➡️  Swagger UI available at: http://localhost:" + os.Getenv("APP_PORT") + "/swagger/index.html")


	// 5️ Middleware
	if os.Getenv("APP_ENV") != "production" {
		app.Use(logger.New())
	}
	app.Use(recover.New())
	app.Use(cors.New())

	// 6️ Root sederhana
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to Alumni API 🚀")
	})

	// 7️7 Register routes (Postgres + Mongo)
	route.SetupPekerjaanMongoRoutes(app, database.MongoDB)
	route.SetupAlumniMongoRoutes(app, database.MongoDB)
	route.SetupRoutes(app, database.DB, database.MongoDB)

	// 8 Tambahkan fitur Upload File
	app.Static("/uploads", "./uploads") // agar file bisa diakses langsung via URL
	uploadRepo := repository.NewFileRepository(database.MongoDB)
	uploadService := service.NewFileService(uploadRepo, "./uploads")
	route.SetupFileRoutes(app, uploadService)

	// 9 Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	go func() {
		log.Printf("Server running on :%s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Server stopped: %v", err)
		}
	}()

	// 10 Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}
}
