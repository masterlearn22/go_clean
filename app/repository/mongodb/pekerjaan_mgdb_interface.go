package repository

import (
	"context"
	"go_clean/app/models/mongodb"
)

type PekerjaanMongoRepositoryInterface interface {
	Create(ctx context.Context, p *models.PekerjaanMongo) (*models.PekerjaanMongo, error)
	FindAll(ctx context.Context) ([]models.PekerjaanMongo, error)
	FindByID(ctx context.Context, id string) (*models.PekerjaanMongo, error)
	FindByAlumniID(ctx context.Context, alumniID int) ([]models.PekerjaanMongo, error)
	Update(ctx context.Context, id string, p *models.PekerjaanMongo) (*models.PekerjaanMongo, error)
	Delete(ctx context.Context, id string) error
}
