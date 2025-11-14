package service

import (
	"fmt"
	"database/sql"
	"go_clean/app/models/postgresql"
	"go_clean/app/repository/postgresql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PekerjaanService struct {
	Repo *repository.PekerjaanRepository
}

// GetAllPekerjaan godoc
// @Summary Ambil semua data pekerjaan
// @Description Mengambil semua pekerjaan dari PostgreSQL (tanpa filter/pagination)
// @Tags Pekerjaan-PostgresSQL
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.PekerjaanAlumni
// @Failure 500 {object} models.ErrorResponse
// @Router /pekerjaan [get]
func (s *PekerjaanService) GetAllPekerjaan(c *fiber.Ctx) error {
	pekerjaan, err := s.Repo.GetAllPekerjaan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan berhasil diambil",
		"data":    pekerjaan,
	})
}

// GetPekerjaanByID godoc
// @Summary Ambil pekerjaan berdasarkan ID
// @Description Mengambil satu pekerjaan berdasarkan ID
// @Tags Pekerjaan-PostgresSQL
// @Security BearerAuth
// @Param id path int true "ID Pekerjaan"
// @Produce json
// @Success 200 {object} models.PekerjaanAlumni
// @Failure 404 {object} models.ErrorResponse
// @Router /pekerjaan/{id} [get]
func (s *PekerjaanService) GetPekerjaanByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID pekerjaan tidak valid",
		})
	}

	pekerjaan, err := s.Repo.GetPekerjaanByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Pekerjaan tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan berhasil diambil",
		"data":    pekerjaan,
	})
}

// GetPekerjaanList godoc
// @Summary Ambil data pekerjaan dengan pagination, search, dan sorting
// @Description Pagination + sorting + search pekerjaan berdasarkan nama_perusahaan atau posisi_jabatan
// @Tags Pekerjaan-PostgresSQL
// @Security BearerAuth
// @Produce json
// @Param search query string false "Cari pekerjaan"
// @Param sortBy query string false "Kolom sorting (nama_perusahaan, posisi_jabatan)"
// @Param order query string false "asc / desc"
// @Param page query int false "Halaman"
// @Param limit query int false "Limit data"
// @Success 200 {object} models.UserResponse[models.PekerjaanAlumni]
// @Failure 500 {object} models.ErrorResponse
// @Router /pekerjaan/list [get]
func (s *PekerjaanService) GetPekerjaanList(c *fiber.Ctx) error {
	sortable := repository.PekerjaanSortable()
	params := getListParams(c, sortable)

	items, err := repository.ListPekerjaanRepo(params.Search, params.SortBy, params.Order, params.Limit, params.Offset)
	if err != nil {
		fmt.Printf("ListPekerjaanRepo error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch pekerjaan"})
	}

	total, err := repository.CountPekerjaanRepo(params.Search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to count pekerjaan"})
	}

	resp := models.UserResponse[models.PekerjaanAlumni]{
		Data: items,
		Meta: models.MetaInfo{
			Page:  params.Page,
			Limit: params.Limit,
			Total: total,
			Pages: (total + params.Limit - 1) / params.Limit,
			SortBy: params.SortBy,
			Order:  params.Order,
			Search: params.Search,
		},
	}
	return c.JSON(resp)
}

// GetPekerjaanByAlumniID godoc
// @Summary Ambil data pekerjaan berdasarkan ID alumni
// @Description Mengambil semua pekerjaan yang dimiliki alumni tertentu
// @Tags Pekerjaan-PostgresSQL
// @Security BearerAuth
// @Param alumni_id path int true "ID Alumni"
// @Produce json
// @Success 200 {array} models.PekerjaanAlumni
// @Failure 404 {object} models.ErrorResponse
// @Router /pekerjaan/alumni/{alumni_id} [get]
func (s *PekerjaanService) GetPekerjaanByAlumniID(c *fiber.Ctx) error {
	alumniID, err := strconv.Atoi(c.Params("alumni_id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID alumni tidak valid",
		})
	}

	pekerjaan, err := s.Repo.GetPekerjaanByAlumniID(alumniID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data pekerjaan: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data pekerjaan untuk alumni berhasil diambil",
		"data":    pekerjaan,
	})
}

// CreatePekerjaan godoc
// @Summary Tambah data pekerjaan baru
// @Description Tambah pekerjaan (hanya bisa diakses Admin)
// @Tags Pekerjaan-PostgresSQL
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.PekerjaanAlumni true "Data pekerjaan baru"
// @Success 201 {object} models.PekerjaanAlumni
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /pekerjaan [post]
func (s *PekerjaanService) CreatePekerjaan(c *fiber.Ctx) error {
	var p models.PekerjaanAlumni
	if err := c.BodyParser(&p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	if p.AlumniID == 0 || p.NamaPerusahaan == "" || p.PosisiJabatan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field alumni_id, nama_perusahaan, dan posisi_jabatan wajib diisi",
		})
	}

	newID, err := s.Repo.CreatePekerjaan(&p)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menambah pekerjaan: " + err.Error(),
		})
	}

	newPekerjaan, _ := s.Repo.GetPekerjaanByID(newID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Pekerjaan berhasil ditambahkan",
		"data":    newPekerjaan,
	})
}

// UpdatePekerjaan godoc
// @Summary Update data pekerjaan
// @Description User hanya boleh update datanya sendiri, Admin boleh update data siapa saja
// @Tags Pekerjaan-PostgresSQL
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Pekerjaan"
// @Param request body models.PekerjaanAlumni true "Data pekerjaan baru"
// @Success 200 {object} models.PekerjaanAlumni
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /pekerjaan/{id} [put]
func (s *PekerjaanService) UpdatePekerjaan(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "ID pekerjaan tidak valid"})
    }

    role := c.Locals("role").(string)
    userID := c.Locals("user_id").(int)

    userRepo := repository.UserRepository{DB: s.Repo.DB}
    user, err := userRepo.GetUserByID(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Gagal mengambil data user"})
    }

    var p models.PekerjaanAlumni
    if err := c.BodyParser(&p); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Request body tidak valid"})
    }

    existing, err := s.Repo.GetPekerjaanByID(id)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Data pekerjaan tidak ditemukan"})
    }

    if role == "user" && existing.AlumniID != *user.AlumniID {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": true, "message": "Tidak punya izin mengubah pekerjaan ini"})
    }

    rows, err := s.Repo.UpdatePekerjaan(id, &p)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Gagal mengupdate pekerjaan"})
    }

    if rows == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "Pekerjaan tidak ditemukan"})
    }

    updated, _ := s.Repo.GetPekerjaanByID(id)
    return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil diupdate", "data": updated})
}

// DeletePekerjaan godoc
// @Summary Soft delete pekerjaan
// @Description Menghapus pekerjaan (soft delete). Admin bisa hapus siapa saja, User hanya boleh data miliknya.
// @Tags Pekerjaan-PostgresSQL
// @Security BearerAuth
// @Param id path int true "ID Pekerjaan"
// @Success 200 {string} string "Pekerjaan berhasil dihapus"
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /pekerjaan/{id} [delete]
func (s *PekerjaanService) DeletePekerjaan(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID pekerjaan tidak valid"})
	}

	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	userRepo := repository.UserRepository{DB: s.Repo.DB}
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data user"})
	}

	existing, err := s.Repo.GetPekerjaanByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Pekerjaan tidak ditemukan"})
	}

	if role == "user" && existing.AlumniID != *user.AlumniID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Tidak punya izin menghapus pekerjaan ini"})
	}

	rows, err := s.Repo.SoftDeletePekerjaan(id, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal menghapus pekerjaan"})
	}

	if rows == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Pekerjaan tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus (soft delete)"})
}

// TrashAllPekerjaan godoc
// @Summary Ambil semua data yang terhapus (soft delete)
// @Description Admin melihat semua data trash, User hanya melihat miliknya
// @Tags Pekerjaan-PostgresSQL
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.PekerjaanAlumni
// @Router /pekerjaan/trash [get]
func (s *PekerjaanService) TrashAllPekerjaan(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)

	userRepo := repository.UserRepository{DB: s.Repo.DB}
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data user"})
	}

	var pekerjaan []models.PekerjaanAlumni
	if role == "admin" {
		pekerjaan, err = s.Repo.TrashAllPekerjaan()
	} else {
		pekerjaan, err = s.Repo.TrashPekerjaanByAlumniID(*user.AlumniID)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data pekerjaan"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Data pekerjaan trash berhasil diambil", "data": pekerjaan})
}

// RestorePekerjaan godoc
// @Summary Restore pekerjaan dari trash
// @Description Mengembalikan pekerjaan (soft delete → active)
// @Tags Pekerjaan-PostgresSQL
// @Security BearerAuth
// @Param id path int true "ID Pekerjaan"
// @Success 200 {string} string "Pekerjaan berhasil direstore"
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /pekerjaan/restore/{id} [put]
func (s *PekerjaanService) RestorePekerjaan(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)
	pekerjaanID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID pekerjaan tidak valid"})
	}

	userRepo := repository.UserRepository{DB: s.Repo.DB}
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data user"})
	}

	existing, err := s.Repo.GetPekerjaanByID(pekerjaanID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Data pekerjaan tidak ditemukan"})
	}

	if role == "user" && existing.AlumniID != *user.AlumniID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Tidak punya izin restore"})
	}

	err = s.Repo.RestorePekerjaanByID(pekerjaanID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal restore pekerjaan"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil di-restore"})
}

// HardDeletePekerjaan godoc
// @Summary Hapus permanen pekerjaan
// @Description Menghapus data secara permanen dari database
// @Tags Pekerjaan-PostgresSQL
// @Security BearerAuth
// @Param id path int true "ID Pekerjaan"
// @Success 200 {string} string "Pekerjaan berhasil dihapus permanen"
// @Failure 403 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /pekerjaan/hard-delete/{id} [delete]
func (s *PekerjaanService) HardDeletePekerjaan(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	userID := c.Locals("user_id").(int)
	pekerjaanID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "ID pekerjaan tidak valid"})
	}

	userRepo := repository.UserRepository{DB: s.Repo.DB}
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal mengambil data user"})
	}

	existing, err := s.Repo.GetPekerjaanByID(pekerjaanID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "Pekerjaan tidak ditemukan"})
	}

	if role == "user" && existing.AlumniID != *user.AlumniID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Tidak punya izin menghapus permanen"})
	}

	err = s.Repo.HardDeletePekerjaanByID(pekerjaanID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": "Gagal menghapus permanen pekerjaan"})
	}

	return c.JSON(fiber.Map{"success": true, "message": "Pekerjaan berhasil dihapus permanen"})
}
