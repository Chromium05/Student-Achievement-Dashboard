package repository

import (
	"student-report/app/model"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func GetAllPekerjaanAlumni(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.PekerjaanAlumni, error) {
	query := fmt.Sprintf(`
		SELECT id , alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_deleted
		FROM pekerjaan_alumni
		WHERE CAST(alumni_id AS TEXT) ILIKE $1 OR posisi_jabatan ILIKE $1 OR bidang_industri ILIKE $1
		OR status_pekerjaan ILIKE $1 OR nama_perusahaan ILIKE $1 AND is_deleted = FALSE
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
		`, sortBy, order)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pekerjaanAlumniList []model.PekerjaanAlumni

	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.Alumni_ID, &pekerjaan.Nama_Perusahaan, &pekerjaan.Posisi_Jabatan, &pekerjaan.Bidang_Industri,
			&pekerjaan.Lokasi_Kerja, &pekerjaan.Gaji_Range,
			&pekerjaan.Mulai_Kerja, &pekerjaan.Selesai_Kerja, &pekerjaan.Status_Pekerjaan, &pekerjaan.Jobdesk,
			&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt, &pekerjaan.IsDeleted,
		)

		if err != nil {
			return nil, err
		}
		pekerjaanAlumniList = append(pekerjaanAlumniList, pekerjaan)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pekerjaanAlumniList, nil
}

func GetPekerjaanByID(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
	// Generate query Postgresql
	row := db.QueryRow(`
		SELECT id , alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE id = $1
	`, id)

	var pekerjaan model.PekerjaanAlumni
	err := row.Scan(
		&pekerjaan.ID, &pekerjaan.Alumni_ID, &pekerjaan.Nama_Perusahaan, &pekerjaan.Posisi_Jabatan, &pekerjaan.Bidang_Industri,
		&pekerjaan.Lokasi_Kerja, &pekerjaan.Gaji_Range,
		&pekerjaan.Mulai_Kerja, &pekerjaan.Selesai_Kerja, &pekerjaan.Status_Pekerjaan, &pekerjaan.Jobdesk,
		&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Pekerjaan tidak ada
		} else {
			return nil, err // Error lainnya
		}
	}

	return &pekerjaan, nil
}

func GetPekerjaanByAlumniID(db *sql.DB, alumniID int) ([]model.PekerjaanAlumni, error) {
	rows, err := db.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE alumni_id = $1
	`, alumniID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pekerjaanAlumniList []model.PekerjaanAlumni
	for rows.Next() {
		var pekerjaan model.PekerjaanAlumni
		err := rows.Scan(
			&pekerjaan.ID, &pekerjaan.Alumni_ID, &pekerjaan.Nama_Perusahaan, &pekerjaan.Posisi_Jabatan, &pekerjaan.Bidang_Industri,
			&pekerjaan.Lokasi_Kerja, &pekerjaan.Gaji_Range,
			&pekerjaan.Mulai_Kerja, &pekerjaan.Selesai_Kerja, &pekerjaan.Status_Pekerjaan, &pekerjaan.Jobdesk,
			&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanAlumniList = append(pekerjaanAlumniList, pekerjaan)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return pekerjaanAlumniList, nil
}

func PostNewPekerjaanAlumni(db *sql.DB, pekerjaan *model.CreatePekerjaan) (*model.PekerjaanAlumni, error) {
	// Validasi input kosong
	if pekerjaan.Alumni_ID == 0 || strings.TrimSpace(pekerjaan.Nama_Perusahaan) == "" ||
		strings.TrimSpace(pekerjaan.Posisi_Jabatan) == "" || strings.TrimSpace(pekerjaan.Bidang_Industri) == "" ||
		strings.TrimSpace(pekerjaan.Lokasi_Kerja) == "" || strings.TrimSpace(pekerjaan.Gaji_Range) == "" ||
		strings.TrimSpace(pekerjaan.Status_Pekerjaan) == "" || strings.TrimSpace(pekerjaan.Jobdesk) == "" {
		return nil, ErrInvalidInput
	}

	// Konversi tanggal selesai kerja: jika kosong, jadikan nil
	var selesaiKerja interface{}
	if strings.TrimSpace(pekerjaan.Selesai_Kerja) == "" {
		selesaiKerja = nil
	} else {
		selesaiKerja = pekerjaan.Selesai_Kerja
	}

	var newID int
	err := db.QueryRow(`
		INSERT INTO pekerjaan_alumni (alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		tanggal_mulai_kerja, tanggal_selesai_kerja,
		status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id
	`, pekerjaan.Alumni_ID, pekerjaan.Nama_Perusahaan, pekerjaan.Posisi_Jabatan, pekerjaan.Bidang_Industri,
		pekerjaan.Lokasi_Kerja, pekerjaan.Gaji_Range, pekerjaan.Mulai_Kerja, selesaiKerja, pekerjaan.Status_Pekerjaan, pekerjaan.Jobdesk).Scan(&newID)

	if err != nil {
		return nil, err
	}

	// Ambil data pekerjaan yang baru ditambahkan
	var newPekerjaan model.PekerjaanAlumni
	row := db.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni WHERE id = $1
	`, newID)

	err = row.Scan(
		&newPekerjaan.ID, &newPekerjaan.Alumni_ID, &newPekerjaan.Nama_Perusahaan, &newPekerjaan.Posisi_Jabatan, &newPekerjaan.Bidang_Industri,
		&newPekerjaan.Lokasi_Kerja, &newPekerjaan.Gaji_Range,
		&newPekerjaan.Mulai_Kerja, &newPekerjaan.Selesai_Kerja, &newPekerjaan.Status_Pekerjaan, &newPekerjaan.Jobdesk,
		&newPekerjaan.CreatedAt, &newPekerjaan.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &newPekerjaan, nil
}

func UpdatePekerjaanAlumni(db *sql.DB, pekerjaan *model.CreatePekerjaan, alumniID int) (*model.PekerjaanAlumni, error) {
	// Validasi input kosong
	if pekerjaan.Alumni_ID == 0 || strings.TrimSpace(pekerjaan.Nama_Perusahaan) == "" ||
		strings.TrimSpace(pekerjaan.Posisi_Jabatan) == "" || strings.TrimSpace(pekerjaan.Bidang_Industri) == "" ||
		strings.TrimSpace(pekerjaan.Lokasi_Kerja) == "" || strings.TrimSpace(pekerjaan.Gaji_Range) == "" ||
		strings.TrimSpace(pekerjaan.Status_Pekerjaan) == "" || strings.TrimSpace(pekerjaan.Jobdesk) == "" {
		return nil, ErrInvalidInput
	}

	// Konversi tanggal selesai kerja: jika kosong, jadikan nil
	var selesaiKerja interface{}
	if strings.TrimSpace(pekerjaan.Selesai_Kerja) == "" {
		selesaiKerja = nil
	} else {
		selesaiKerja = pekerjaan.Selesai_Kerja
	}

	_, err := db.Exec(`
		UPDATE pekerjaan_alumni
		SET alumni_id = $1, nama_perusahaan = $2, posisi_jabatan = $3, bidang_industri = $4, lokasi_kerja = $5, 
		gaji_range = $6, tanggal_mulai_kerja = $7, tanggal_selesai_kerja = $8, status_pekerjaan = $9,
		deskripsi_pekerjaan = $10, updated_at = $11
		WHERE id = $12
	`, pekerjaan.Alumni_ID, pekerjaan.Nama_Perusahaan, pekerjaan.Posisi_Jabatan, pekerjaan.Bidang_Industri,
		pekerjaan.Lokasi_Kerja, pekerjaan.Gaji_Range, pekerjaan.Mulai_Kerja, selesaiKerja, pekerjaan.Status_Pekerjaan,
		pekerjaan.Jobdesk, time.Now(), alumniID)

	if err != nil {
		return nil, err
	}

	// Ambil data pekerjaan yang sudah diupdate
	var updatedPekerjaan model.PekerjaanAlumni
	row := db.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni WHERE id = $1
	`, alumniID)

	err = row.Scan(
		&updatedPekerjaan.ID, &updatedPekerjaan.Alumni_ID, &updatedPekerjaan.Nama_Perusahaan, &updatedPekerjaan.Posisi_Jabatan, &updatedPekerjaan.Bidang_Industri,
		&updatedPekerjaan.Lokasi_Kerja, &updatedPekerjaan.Gaji_Range,
		&updatedPekerjaan.Mulai_Kerja, &updatedPekerjaan.Selesai_Kerja, &updatedPekerjaan.Status_Pekerjaan, &updatedPekerjaan.Jobdesk,
		&updatedPekerjaan.CreatedAt, &updatedPekerjaan.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updatedPekerjaan, nil
}

func DeletePekerjaanAlumni(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
	// Ambil data sebelum dihapus
	var pekerjaan model.PekerjaanAlumni
	row := db.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni WHERE id = $1
	`, id)

	err := row.Scan(
		&pekerjaan.ID, &pekerjaan.Alumni_ID, &pekerjaan.Nama_Perusahaan, &pekerjaan.Posisi_Jabatan, &pekerjaan.Bidang_Industri,
		&pekerjaan.Lokasi_Kerja, &pekerjaan.Gaji_Range,
		&pekerjaan.Mulai_Kerja, &pekerjaan.Selesai_Kerja, &pekerjaan.Status_Pekerjaan, &pekerjaan.Jobdesk,
		&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Query SQL untuk menghapus data
	_, err = db.Exec(`
		DELETE FROM pekerjaan_alumni WHERE id = $1 AND is_deleted = TRUE
	`, id)

	if err != nil {
		return nil, err
	}

	return &pekerjaan, nil
}

func SoftDeletePekerjaanAlumni(db *sql.DB, id int) (*model.SoftDeletePekerjaanAlumni, error) {
	var pekerjaan model.SoftDeletePekerjaanAlumni
	row := db.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at, is_deleted
		FROM pekerjaan_alumni WHERE id = $1
	`, id)

	err := row.Scan(
		&pekerjaan.ID, &pekerjaan.Alumni_ID, &pekerjaan.Nama_Perusahaan, &pekerjaan.Posisi_Jabatan, &pekerjaan.Bidang_Industri,
		&pekerjaan.Lokasi_Kerja, &pekerjaan.Gaji_Range,
		&pekerjaan.Mulai_Kerja, &pekerjaan.Selesai_Kerja, &pekerjaan.Status_Pekerjaan, &pekerjaan.Jobdesk,
		&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt, &pekerjaan.IsDeleted,
	)
	if err != nil {
		return nil, err
	}

	// Query SQL untuk menghapus data
	_, err = db.Exec(`
		UPDATE pekerjaan_alumni
		SET is_deleted = TRUE
		WHERE id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	return &pekerjaan, nil
}
