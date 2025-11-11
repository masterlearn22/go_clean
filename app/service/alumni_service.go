package service

import (
	"database/sql"
	"fmt"
	"go_clean/app/models"
	"go_clean/app/repository"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AlumniService struct {
	Repo *repository.AlumniRepository
}

// GetAllAlumni godoc
// @Summary Ambil semua data alumni (PostgreSQL)
// @Description Menampilkan semua alumni dari database PostgreSQL
// @Tags Alumni-PostgresSQL
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Alumni
// @Failure 500 {object} models.ErrorResponse
// @Router /alumni [get]
func (s *AlumniService) GetAllAlumni(c *fiber.Ctx) error {
	alumni, err := s.Repo.GetAllAlumni()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil",
		"data":    alumni,
	})
}

// GetAlumniList godoc
// @Summary Dapatkan alumni dengan pagination, sorting & searching
// @Description Pagination + sorting + search alumni berdasarkan nama atau NIM
// @Tags Alumni-PostgresSQL
// @Security BearerAuth
// @Produce json
// @Param search query string false "cari nama atau nim"
// @Param sortBy query string false "sort berdasarkan kolom (nim,nama,jurusan,angkatan)"
// @Param order query string false "asc atau desc"
// @Param page query int false "Halaman"
// @Param limit query int false "Limit data"
// @Success 200 {object} models.UserResponse[models.Alumni]
// @Failure 500 {object} models.ErrorResponse
// @Router /alumni-pag [get]
func (s *AlumniService) GetAlumniList(c *fiber.Ctx) error {
	sortable := make(map[string]bool)
	for _, v := range repository.AlumniSortable() {
		sortable[v] = true
	}
	params := getListParams(c, sortable)
	items, err := repository.ListAlumniRepo(params.Search, params.SortBy, params.Order, params.Limit, params.Offset)
	if err != nil {
		fmt.Printf("ListAlumniRepo error: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch alumni",
		})
	}

	total, err := repository.CountAlumniRepo(params.Search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to count alumni"})
	}

	resp := models.UserResponse[models.Alumni]{
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

// GetAlumniByID godoc
// @Summary Ambil alumni berdasarkan ID
// @Description Mengambil satu alumni dari PostgreSQL berdasarkan ID
// @Tags Alumni-PostgresSQL
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID Alumni"
// @Success 200 {object} models.Alumni
// @Failure 404 {object} models.ErrorResponse
// @Router /alumni/{id} [get]
func (s *AlumniService) GetAlumniByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	alumni, err := s.Repo.GetAlumniByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"success": false,
				"message": "Alumni tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni: " + err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil",
		"data":    alumni,
	})
}

// GetAlumniByAngkatan godoc
// @Summary Ambil alumni berdasarkan angkatan
// @Description Mengambil jumlah alumni berdasarkan angkatan tertentu
// @Tags Alumni-PostgresSQL
// @Security BearerAuth
// @Produce json
// @Param angkatan path int true "Tahun angkatan"
// @Success 200 {object} models.AlumniAngkatan
// @Failure 404 {object} models.ErrorResponse
// @Router /alumni/angkatan/{angkatan} [get]
func (s *AlumniService) GetAlumniByAngkatan(c *fiber.Ctx) error {
	angkatan, err := strconv.Atoi(c.Params("angkatan"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Angkatan tidak valid",
		})
	}

	result, err := s.Repo.GetAlumniByAngkatan(angkatan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni berhasil diambil",
		"data":    result,
	})
}

// GetAlumniAndPekerjaan godoc
// @Summary Ambil alumni beserta data pekerjaan
// @Description Join data alumni & pekerjaan dari PostgreSQL
// @Tags Alumni-PostgresSQL
// @Security BearerAuth
// @Produce json
// @Param nim path int true "NIM Alumni"
// @Success 200 {object} models.AlumniPekerjaan
// @Failure 404 {object} models.ErrorResponse
// @Router /alumni/detail/{nim} [get]
func (s *AlumniService) GetAlumniAndPekerjaan(c *fiber.Ctx) error {
	idStr := c.Params("nim")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID harus berupa angka",
		})
	}

	result, err := s.Repo.GetAlumniAndPekerjaan(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data alumni dan pekerjaan: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Data alumni dan pekerjaan berhasil diambil",
		"data":    result,
	})
}

// CreateAlumni godoc
// @Summary Tambah alumni baru
// @Description Insert data alumni ke database PostgreSQL
// @Tags Alumni-PostgresSQL
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body models.Alumni true "Data alumni baru"
// @Success 201 {object} models.Alumni
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /alumni [post]
func (s *AlumniService) CreateAlumni(c *fiber.Ctx) error {
	var alumni models.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	if alumni.NIM == "" || alumni.Nama == "" || alumni.Jurusan == "" || alumni.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Field NIM, Nama, Jurusan, dan Email wajib diisi",
		})
	}

	newID, err := s.Repo.CreateAlumni(&alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menambah alumni: " + err.Error(),
		})
	}

	newAlumni, _ := s.Repo.GetAlumniByID(newID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil ditambahkan",
		"data":    newAlumni,
	})
}

// UpdateAlumni godoc
// @Summary Update alumni
// @Description Mengubah data alumni berdasarkan ID
// @Tags Alumni-PostgresSQL
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID Alumni"
// @Param request body models.Alumni true "Data alumni baru"
// @Success 200 {object} models.Alumni
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /alumni/{id} [put]
func (s *AlumniService) UpdateAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	var alumni models.Alumni
	if err := c.BodyParser(&alumni); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Request body tidak valid",
		})
	}

	rowsAffected, err := s.Repo.UpdateAlumni(id, &alumni)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengupdate alumni: " + err.Error(),
		})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Alumni tidak ditemukan untuk diupdate",
		})
	}

	updatedAlumni, _ := s.Repo.GetAlumniByID(id)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil diupdate",
		"data":    updatedAlumni,
	})
}

// DeleteAlumni godoc
// @Summary Hapus alumni
// @Description Menghapus data alumni berdasarkan ID
// @Tags Alumni-PostgresSQL
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID Alumni"
// @Success 204 "No Content"
// @Failure 404 {object} models.ErrorResponse
// @Router /alumni/{id} [delete]
func (s *AlumniService) DeleteAlumni(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "ID tidak valid",
		})
	}

	rowsAffected, err := s.Repo.DeleteAlumni(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus alumni: " + err.Error(),
		})
	}
	if rowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Alumni tidak ditemukan untuk dihapus",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Alumni berhasil dihapus",
	})
}
