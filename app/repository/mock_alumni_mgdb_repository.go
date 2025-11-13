package repository

import (
	"context"
	"errors"
	"go_clean/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAlumniMongoRepository struct {
	Data map[string]*models.AlumniMongo
}

func NewMockAlumniMongoRepository() *MockAlumniMongoRepository {
	return &MockAlumniMongoRepository{
		Data: make(map[string]*models.AlumniMongo),
	}
}

func (m *MockAlumniMongoRepository) Create(ctx context.Context, data *models.AlumniMongo) (*models.AlumniMongo, error) {
	if data.Nama == "" {
		return nil, errors.New("nama tidak boleh kosong")
	}

	data.ID = primitive.NewObjectID()
	m.Data[data.ID.Hex()] = data
	return data, nil
}

func (m *MockAlumniMongoRepository) FindAll(ctx context.Context) ([]models.AlumniMongo, error) {
	var list []models.AlumniMongo
	for _, v := range m.Data {
		list = append(list, *v)
	}
	return list, nil
}

func (m *MockAlumniMongoRepository) FindByID(ctx context.Context, id string) (*models.AlumniMongo, error) {
	if val, ok := m.Data[id]; ok {
		return val, nil
	}
	return nil, errors.New("data tidak ditemukan")
}

func (m *MockAlumniMongoRepository) Update(ctx context.Context, id string, data *models.AlumniMongo) (*models.AlumniMongo, error) {
	if _, ok := m.Data[id]; !ok {
		return nil, errors.New("data tidak ditemukan")
	}
	m.Data[id] = data
	return data, nil
}

func (m *MockAlumniMongoRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.Data[id]; !ok {
		return errors.New("data tidak ditemukan")
	}
	delete(m.Data, id)
	return nil
}
