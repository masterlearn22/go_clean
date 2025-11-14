package service

import (
	"context"
	"go_clean/app/models/mongodb"
	"go_clean/app/repository/mongodb"
	"time"
)

type AlumniMongoService struct {
	repo repository.AlumniMongoRepositoryInterface
}

func NewAlumniMongoService(repo repository.AlumniMongoRepositoryInterface) *AlumniMongoService {
	return &AlumniMongoService{repo: repo}
}

// Create godoc
// @Summary Tambahkan alumni baru
// @Description Menambahkan data alumni ke database MongoDB
// @Tags Alumni-Mongo
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.AlumniMongo true "Data alumni"
// @Success 201 {object} models.AlumniMongo
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /alumni-mongo [post]
func (s *AlumniMongoService) Create(ctx context.Context, data *models.AlumniMongo) (*models.AlumniMongo, error) {
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	return s.repo.Create(ctx, data)
}

// GetAllAlumni godoc
// @Summary Dapatkan semua alumni
// @Description Mengambil daftar semua alumni dari database
// @Tags Alumni-Mongo
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.AlumniMongo
// @Failure 500 {object} models.ErrorResponse
// @Router  /alumni-mongo [get]
func (s *AlumniMongoService) GetAll(ctx context.Context) ([]models.AlumniMongo, error) {
	return s.repo.FindAll(ctx)
}

// GetByID godoc
// @Summary Dapatkan alumni berdasarkan ID
// @Description Mengambil satu alumni berdasarkan ID MongoDB
// @Tags Alumni-Mongo
// @Security BearerAuth
// @Produce json
// @Param id path string true "Alumni ID"
// @Success 200 {object} models.AlumniMongo
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /alumni-mongo/{id} [get]
func (s *AlumniMongoService) GetByID(ctx context.Context, id string) (*models.AlumniMongo, error) {
	return s.repo.FindByID(ctx, id)
}

// Update godoc
// @Summary Update data alumni
// @Description Mengubah data alumni berdasarkan ID
// @Tags Alumni-Mongo
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Alumni ID"
// @Param request body models.AlumniMongo true "Data alumni baru"
// @Success 200 {object} models.AlumniMongo
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /alumni-mongo/{id} [put]
func (s *AlumniMongoService) Update(ctx context.Context, id string, data *models.AlumniMongo) (*models.AlumniMongo, error) {
	data.UpdatedAt = time.Now()
	return s.repo.Update(ctx, id, data)
}

// Delete godoc
// @Summary Hapus alumni
// @Description Menghapus data alumni dari database berdasarkan ID
// @Tags Alumni-Mongo
// @Security BearerAuth
// @Produce json
// @Param id path string true "Alumni ID"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /alumni-mongo/{id} [delete]
func (s *AlumniMongoService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
