package repository

import (
	"errors"
	"go_clean/app/models/mongodb"
)

type MockFileRepository struct {
	Data  map[string]*models.File
	ForceError bool // untuk memaksa error dari DB
}

func NewMockFileRepository() *MockFileRepository {
	return &MockFileRepository{
		Data: make(map[string]*models.File),
	}
}

func (m *MockFileRepository) Create(file *models.File) error {
	if m.ForceError {
		return errors.New("database error")
	}

	m.Data[file.ID.Hex()] = file
	return nil
}

func (m *MockFileRepository) FindAll() ([]models.File, error) {
	var files []models.File
	for _, f := range m.Data {
		files = append(files, *f)
	}
	return files, nil
}

func (m *MockFileRepository) FindByID(id string) (*models.File, error) {
	if f, ok := m.Data[id]; ok {
		return f, nil
	}
	return nil, errors.New("not found")
}

func (m *MockFileRepository) Delete(id string) error {
	delete(m.Data, id)
	return nil
}
