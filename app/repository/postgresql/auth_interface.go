package repository

import "go_clean/app/models/postgresql"

type IUserRepository interface {
	GetByUsernameOrEmail(identifier string) (*models.User, string, error)
	ExistsByUsernameOrEmail(username, email string) (bool, error)
	Create(username, email, passwordHash, role string) (*models.User, error)
	GetUsersRepo(search, sortBy, order string, limit, offset int) ([]models.User, error)
	CountUsersRepo(search string) (int, error)
}
