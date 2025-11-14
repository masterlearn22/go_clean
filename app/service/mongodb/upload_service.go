package service

import (
    "fmt"
    "os"
    "path/filepath"
    "go_clean/app/models/mongodb"
    "go_clean/app/repository/mongodb"

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

// UploadFile godoc
// @Summary Upload file (PDF / Image)
// @Description Mengupload file ke server dan menyimpan metadata ke database (MongoDB)
// @Tags FileUpload
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File yang akan diupload"
// @Success 201 {object} models.FileResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /files/upload [post]
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
        "data":    fileModel,
    })
}


// GetAllFiles godoc
// @Summary Mendapatkan semua file yang sudah diupload
// @Tags FileUpload
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.FileResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /files/ [get]
func (s *fileService) GetAllFiles(c *fiber.Ctx) error {
    files, err := s.repo.FindAll()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "success": false,
            "message": "Failed to get files",
            "error":   err.Error(),
        })
    }

    return c.JSON(fiber.Map{
        "success": true,
        "message": "Files retrieved successfully",
        "data":    files,
    })
}


// GetFileByID godoc
// @Summary Mendapatkan file berdasarkan ID
// @Tags FileUpload
// @Security BearerAuth
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} models.FileResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /files/{id} [get]
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
        "data":    file,
    }) 
}
// DeleteFile godoc
// @Summary Menghapus file berdasarkan ID
// @Security BearerAuth
// @Tags FileUpload
// @Security BearerAuth
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /files/{id} [delete]
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