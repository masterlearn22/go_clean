package repository

import (
	"context"
	"go_clean/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoRepository struct {
	Col *mongo.Collection
}

func NewUserMongoRepository(db *mongo.Database) *UserMongoRepository {
	return &UserMongoRepository{
		Col: db.Collection("users"),
	}
}

func (r *UserMongoRepository) FindByUsernameOrEmail(identifier string) (*models.LoginMongo, error) {
	var user models.LoginMongo
	err := r.Col.FindOne(
		context.TODO(),
		bson.M{
			"$or": []bson.M{
				{"username": identifier},
				{"email": identifier},
			},
		},
	).Decode(&user)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserMongoRepository) CreateUser(u *models.LoginMongo) (*models.LoginMongo, error) {
	_, err := r.Col.InsertOne(context.TODO(), u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
