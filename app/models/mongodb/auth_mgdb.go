package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoginMongo struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username"`
	Email     string             `bson:"email" json:"email"`
	PasswordHash  string             `bson:"password_hash" json:"-"` // hash
	Role      string             `bson:"role" json:"role"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}
