package repository

import "go_clean/app/models/mongodb"

type AuthRepositoryInterface interface {
	FindByUsernameOrEmail(identifier string) (*models.LoginMongo, error)
}
