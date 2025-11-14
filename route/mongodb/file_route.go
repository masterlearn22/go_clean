package route

import (
    "go_clean/app/service/mongodb"
    "github.com/gofiber/fiber/v2"
)

func SetupFileRoutes(app *fiber.App, service service.FileService) {
    api := app.Group("/api")

    files := api.Group("/files")
    files.Post("/upload", service.UploadFile)
    files.Get("/", service.GetAllFiles)
    files.Get("/:id", service.GetFileByID)
    files.Delete("/:id", service.DeleteFile)
}
