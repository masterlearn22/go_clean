package service

import (
    "fmt"
    "os"
    "path/filepath"
    "go_clean/app/models"
    "go_clean/app/repository"

    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
)

type FileService interface {
    UploadFile(c *fiber.Ctx) error
    GetAllFiles(c *fiber.Ctx) error
    GetFileByID(c *fiber.Ctx) error
    DeleteFile(c *fiber.Ctx) error
}

type fileService struct {
    repo       repository.FileRepository
    uploadPath string
}

func NewFileService(repo repository.FileRepository, uploadPath string) FileService {
    return &fileService{
        repo:       repo,
        uploadPath: uploadPath,
    }
}

func (s *fileService) UploadFile(c *fiber.Ctx) error {
    fileHeader, err := c.FormFile("file")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "No file uploaded",
            "error":   err.Error(),
        })
    }

    if fileHeader.Size > 10*1024*1024 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "File size exceeds 10MB",
        })
    }

    allowedTypes := map[string]bool{
        "image/jpeg":      true,
        "image/png":       true,
        "image/jpg":       true,
        "application/pdf": true,
    }

    contentType := fileHeader.Header.Get("Content-Type")
    if !allowedTypes[contentType] {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "success": false,
            "message": "File type not allowed",
        })
    }

    ext := filepath.Ext(fileHeader.Filename)
    newFileName := uuid.New().String() + ext
    filePath := filepath.Join(s.uploadPath, newFileName)

    if err := os.MkdirAll(s.uploadPath, os.ModePerm); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to create upload directory",
            "error":   err.Error(),
        })
    }

    if err := c.SaveFile(fileHeader, filePath); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to save file",
            "error":   err.Error(),
        })
    }

    fileModel := &models.File{
        FileName:     newFileName,
        OriginalName: fileHeader.Filename,
        FilePath:     filePath,
        FileSize:     fileHeader.Size,
        FileType:     contentType,
    }

    if err := s.repo.Create(fileModel); err != nil {
        os.Remove(filePath)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to save file metadata",
            "error":   err.Error(),
        })
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "success": true,
        "message": "File uploaded successfully",
        "data":    s.toFileResponse(fileModel),
    })
}

func (s *fileService) GetAllFiles(c *fiber.Ctx) error {
    files, err := s.repo.FindAll()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to get files",
            "error":   err.Error(),
        })
    }

    var responses []models.FileResponse
    for _, file := range files {
        responses = append(responses, *s.toFileResponse(&file))
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Files retrieved successfully",
        "data":    responses,
    })
}

func (s *fileService) GetFileByID(c *fiber.Ctx) error {
    id := c.Params("id")

    file, err := s.repo.FindByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "File not found",
            "error":   err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "File retrieved successfully",
        "data":    s.toFileResponse(file),
    })
}

func (s *fileService) DeleteFile(c *fiber.Ctx) error {
    id := c.Params("id")

    file, err := s.repo.FindByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "success": false,
            "message": "File not found",
            "error":   err.Error(),
        })
    }

    if err := os.Remove(file.FilePath); err != nil {
        fmt.Println("Warning: Failed to delete file:", err)
    }

    if err := s.repo.Delete(id); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to delete file",
            "error":   err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "File deleted successfully",
    })
}

func (s *fileService) toFileResponse(file *models.File) *models.FileResponse {
    return &models.FileResponse{
        ID:           file.ID.Hex(),
        FileName:     file.FileName,
        OriginalName: file.OriginalName,
        FilePath:     file.FilePath,
        FileSize:     file.FileSize,
        FileType:     file.FileType,
        UploadedAt:   file.UploadedAt,
    }
}
