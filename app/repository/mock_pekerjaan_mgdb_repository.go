package repository

import (
	"context"
	"errors"
	"go_clean/app/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockPekerjaanMongoRepository struct {
	Data map[string]*models.PekerjaanMongo
}

func NewMockPekerjaanMongoRepository() *MockPekerjaanMongoRepository {
	return &MockPekerjaanMongoRepository{
		Data: make(map[string]*models.PekerjaanMongo),
	}
}

func (m *MockPekerjaanMongoRepository) Create(ctx context.Context, p *models.PekerjaanMongo) (*models.PekerjaanMongo, error) {
	if p.NamaPerusahaan == "" {
		return nil, errors.New("nama pekerjaan tidak boleh kosong")
	}
	p.ID = primitive.NewObjectID()
	m.Data[p.ID.Hex()] = p
	return p, nil
}

func (m *MockPekerjaanMongoRepository) FindAll(ctx context.Context) ([]models.PekerjaanMongo, error) {
	var list []models.PekerjaanMongo
	for _, v := range m.Data {
		list = append(list, *v)
	}
	return list, nil
}

func (m *MockPekerjaanMongoRepository) FindByID(ctx context.Context, id string) (*models.PekerjaanMongo, error) {
	if v, ok := m.Data[id]; ok {
		return v, nil
	}
	return nil, errors.New("data tidak ditemukan")
}

func (m *MockPekerjaanMongoRepository) FindByAlumniID(ctx context.Context, alumniID int) ([]models.PekerjaanMongo, error) {
	var list []models.PekerjaanMongo
	for _, v := range m.Data {
		if v.AlumniID == alumniID {
			list = append(list, *v)
		}
	}
	return list, nil
}

func (m *MockPekerjaanMongoRepository) Update(ctx context.Context, id string, p *models.PekerjaanMongo) (*models.PekerjaanMongo, error) {
	if _, ok := m.Data[id]; !ok {
		return nil, errors.New("data tidak ditemukan")
	}
	m.Data[id] = p
	return p, nil
}

func (m *MockPekerjaanMongoRepository) Delete(ctx context.Context, id string) error {
	if _, ok := m.Data[id]; !ok {
		return errors.New("data tidak ditemukan")
	}
	delete(m.Data, id)
	return nil
}
