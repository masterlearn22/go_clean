package service

import (
	"context"
	"testing"
	"time"

	"go_clean/app/models/mongodb"
	"go_clean/app/repository/mongodb"
)

func TestCreatePekerjaan(t *testing.T) {
	mockRepo := repository.NewMockPekerjaanMongoRepository()
	svc := NewPekerjaanMongoService(mockRepo)
	ctx := context.Background()

	p := &models.PekerjaanMongo{
		AlumniID:      123,
		NamaPerusahaan: "PT. Tech Indonesia",
		PosisiJabatan:    "BackEnd Dev",
	}

	result, err := svc.Create(ctx, p)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.NamaPerusahaan != "PT. Tech Indonesia" {
		t.Errorf("expected PT. Tech Indonesia, got %v", result.NamaPerusahaan)
	}
}

func TestGetByIDPekerjaan(t *testing.T) {
	mockRepo := repository.NewMockPekerjaanMongoRepository()
	svc := NewPekerjaanMongoService(mockRepo)
	ctx := context.Background()

	p := &models.PekerjaanMongo{NamaPerusahaan: "Engineer"}
	created, _ := svc.Create(ctx, p)

	result, err := svc.GetByID(ctx, created.ID.Hex())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.NamaPerusahaan != "Engineer" {
		t.Errorf("expected Engineer, got %v", result.NamaPerusahaan)
	}
}

func TestGetByAlumniID(t *testing.T) {
	mockRepo := repository.NewMockPekerjaanMongoRepository()
	svc := NewPekerjaanMongoService(mockRepo)
	ctx := context.Background()

	p1 := &models.PekerjaanMongo{AlumniID: 10, NamaPerusahaan: "DevOps"}
	p2 := &models.PekerjaanMongo{AlumniID: 10, NamaPerusahaan: "Data Analyst"}

	svc.Create(ctx, p1)
	svc.Create(ctx, p2)

	list, err := svc.GetByAlumniID(ctx, 10)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(list) != 2 {
		t.Errorf("expected 2 results, got %d", len(list))
	}
}

func TestUpdatePekerjaan(t *testing.T) {
	mockRepo := repository.NewMockPekerjaanMongoRepository()
	svc := NewPekerjaanMongoService(mockRepo)
	ctx := context.Background()

	p := &models.PekerjaanMongo{NamaPerusahaan: "Old Job"}
	created, _ := svc.Create(ctx, p)

	updated := &models.PekerjaanMongo{
		NamaPerusahaan: "New Job",
		UpdatedAt:     time.Now(),
	}

	result, err := svc.Update(ctx, created.ID.Hex(), updated)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.NamaPerusahaan != "New Job" {
		t.Errorf("expected New Job, got %v", result.NamaPerusahaan)
	}
}

func TestDeletePekerjaan(t *testing.T) {
	mockRepo := repository.NewMockPekerjaanMongoRepository()
	svc := NewPekerjaanMongoService(mockRepo)
	ctx := context.Background()

	p := &models.PekerjaanMongo{NamaPerusahaan: "Delete Me"}
	created, _ := svc.Create(ctx, p)

	err := svc.Delete(ctx, created.ID.Hex())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = svc.GetByID(ctx, created.ID.Hex())
	if err == nil {
		t.Errorf("expected error after delete, got nil")
	}
}
