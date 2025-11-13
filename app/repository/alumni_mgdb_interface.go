package repository

import (
	"context"
	"go_clean/app/models"
)

type AlumniMongoRepositoryInterface interface {
	Create(ctx context.Context, data *models.AlumniMongo) (*models.AlumniMongo, error)
	FindAll(ctx context.Context) ([]models.AlumniMongo, error)
	FindByID(ctx context.Context, id string) (*models.AlumniMongo, error)
	Update(ctx context.Context, id string, data *models.AlumniMongo) (*models.AlumniMongo, error)
	Delete(ctx context.Context, id string) error
}
