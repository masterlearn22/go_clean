package service

import (
	"context"
	"go_clean/app/models"
	"go_clean/app/repository"
	"time"
)

type PekerjaanMongoService struct {
	repo repository.PekerjaanMongoRepositoryInterface
}

func NewPekerjaanMongoService(repo repository.PekerjaanMongoRepositoryInterface) *PekerjaanMongoService {
	return &PekerjaanMongoService{repo: repo}
}

// Create godoc
// @Summary Tambahkan data pekerjaan baru
// @Description Menambahkan data pekerjaan ke MongoDB
// @Tags Pekerjaan-Mongo
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.PekerjaanMongo true "Data pekerjaan"
// @Success 201 {object} models.PekerjaanMongo
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /pekerjaan-mongo [post]
func (s *PekerjaanMongoService) Create(ctx context.Context, p *models.PekerjaanMongo) (*models.PekerjaanMongo, error) {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return s.repo.Create(ctx, p)
}

// GetAll godoc
// @Summary Dapatkan semua data pekerjaan
// @Description Mengambil semua data pekerjaan di MongoDB
// @Tags Pekerjaan-Mongo
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.PekerjaanMongo
// @Failure 500 {object} models.ErrorResponse
// @Router /pekerjaan-mongo [get]
func (s *PekerjaanMongoService) GetAll(ctx context.Context) ([]models.PekerjaanMongo, error) {
	return s.repo.FindAll(ctx)
}

// GetByID godoc
// @Summary Dapatkan pekerjaan berdasarkan ID
// @Description Mengambil satu pekerjaan berdasarkan ObjectID MongoDB
// @Tags Pekerjaan-Mongo
// @Security BearerAuth
// @Produce json
// @Param id path string true "Pekerjaan ID"
// @Success 200 {object} models.PekerjaanMongo
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /pekerjaan-mongo/{id} [get]
func (s *PekerjaanMongoService) GetByID(ctx context.Context, id string) (*models.PekerjaanMongo, error) {
	return s.repo.FindByID(ctx, id)
}

// GetByAlumniID godoc
// @Summary Dapatkan daftar pekerjaan berdasarkan AlumniID
// @Description Mengambil semua pekerjaan berdasarkan AlumniID
// @Tags Pekerjaan-Mongo
// @Security BearerAuth
// @Produce json
// @Param alumni_id path int true "ID Alumni"
// @Success 200 {array} models.PekerjaanMongo
// @Failure 500 {object} models.ErrorResponse
// @Router /pekerjaan-mongo/alumni/{alumni_id} [get]
func (s *PekerjaanMongoService) GetByAlumniID(ctx context.Context, alumniID int) ([]models.PekerjaanMongo, error) {
	return s.repo.FindByAlumniID(ctx, alumniID)
}

// Update godoc
// @Summary Update data pekerjaan
// @Description Mengubah data pekerjaan berdasarkan ID
// @Tags Pekerjaan-Mongo
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Pekerjaan ID"
// @Param request body models.PekerjaanMongo true "Data pekerjaan baru"
// @Success 200 {object} models.PekerjaanMongo
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /pekerjaan-mongo/{id} [put]
func (s *PekerjaanMongoService) Update(ctx context.Context, id string, p *models.PekerjaanMongo) (*models.PekerjaanMongo, error) {
	p.UpdatedAt = time.Now()
	return s.repo.Update(ctx, id, p)
}

// Delete godoc
// @Summary Hapus pekerjaan
// @Description Menghapus data pekerjaan dari database berdasarkan ID
// @Tags Pekerjaan-Mongo
// @Security BearerAuth
// @Produce json
// @Param id path string true "Pekerjaan ID"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /pekerjaan-mongo/{id} [delete]
func (s *PekerjaanMongoService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
