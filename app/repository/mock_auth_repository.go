package repository

import (
	"go_clean/app/models"
)

type MockUserRepository struct {
	User        *models.User
	Hash        string
	Err         error
	Exists      bool
	Users       []models.User
	Count       int
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) GetByUsernameOrEmail(identifier string) (*models.User, string, error) {
	if m.Err != nil {
		return nil, "", m.Err
	}
	return m.User, m.Hash, nil
}

func (m *MockUserRepository) ExistsByUsernameOrEmail(username, email string) (bool, error) {
	return m.Exists, m.Err
}

func (m *MockUserRepository) Create(username, email, passwordHash, role string) (*models.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return &models.User{
		ID:       99,
		Username: username,
		Email:    email,
		Role:     role,
	}, nil
}

func (m *MockUserRepository) GetUsersRepo(search, sortBy, order string, limit, offset int) ([]models.User, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Users, nil
}

func (m *MockUserRepository) CountUsersRepo(search string) (int, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	return m.Count, nil
}
