package repository

import (
    "errors"
    "go_clean/app/models/mongodb"
)

type MockUserMongoRepository struct {
    Users map[string]*models.LoginMongo
}

func NewMockUserMongoRepository() *MockUserMongoRepository {
    return &MockUserMongoRepository{
        Users: make(map[string]*models.LoginMongo),
    }
}

func (m *MockUserMongoRepository) FindByUsernameOrEmail(identifier string) (*models.LoginMongo, error) {
    user, ok := m.Users[identifier]
    if !ok {
        return nil, errors.New("not found")
    }
    return user, nil
}

func (m *MockUserMongoRepository) InsertUser(u *models.LoginMongo) {
    m.Users[u.Username] = u
    m.Users[u.Email] = u
}
