package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlumniMongo struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty" example:"123456789abcdefghijklmnopqrstuvwxyz"`
	AlumniID    int                `bson:"alumni_id" json:"alumni_id" example:"1"`
	NIM         string             `bson:"nim" json:"nim" example:"434231048"`
	Nama        string             `bson:"nama" json:"nama" example:"Surya Dwi Satria"`
	Jurusan     string             `bson:"jurusan" json:"jurusan" example:"Teknik Informatika"`
	Angkatan    int                `bson:"angkatan" json:"angkatan" example:"2023"`
	TahunLulus  int                `bson:"tahun_lulus" json:"tahun_lulus" example:"2027"`
	Email       string             `bson:"email" json:"email" example:"surya@example.com"`
	NoTelp      string             `bson:"no_telepon" json:"no_telepon" example:"081234567890"`
	Alamat      string             `bson:"alamat" json:"alamat" example:"Surabaya"`
	TempatKerja string             `bson:"tempat_kerja,omitempty" json:"tempat_kerja,omitempty" example:"PT Tech Indonesia"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at" example:"2025-11-07T11:20:00Z"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at" example:"2025-11-07T11:20:00Z"`
}

