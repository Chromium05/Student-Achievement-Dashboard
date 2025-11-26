package repository

import (
	"student-report/app/model"
	"database/sql"
	// "fmt"
	// "strings"
	// "time"
)

func GetAllTrash(db *sql.DB) ([]model.PekerjaanAlumni, error) {
	rows, err := db.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, gaji_range,
		tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni
		WHERE is_deleted = TRUE
	`)
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

// Hanya untuk validasi ID di auth
func GetTrashByID(db *sql.DB, alumniID int) ([]model.PekerjaanAlumni, error) {
	rows, err := db.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
		lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
		status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
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
			&pekerjaan.ID, &pekerjaan.Alumni_ID, &pekerjaan.Nama_Perusahaan, 
			&pekerjaan.Posisi_Jabatan, &pekerjaan.Bidang_Industri, &pekerjaan.Lokasi_Kerja, 
			&pekerjaan.Gaji_Range, &pekerjaan.Mulai_Kerja, &pekerjaan.Selesai_Kerja, 
			&pekerjaan.Status_Pekerjaan, &pekerjaan.Jobdesk,
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

func RestoreData(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
    var pekerjaan model.PekerjaanAlumni
    row := db.QueryRow(`
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
		lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
		status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
        FROM pekerjaan_alumni WHERE id = $1
    `, id)

    err := row.Scan(
        &pekerjaan.ID, &pekerjaan.Alumni_ID, &pekerjaan.Nama_Perusahaan, 
		&pekerjaan.Posisi_Jabatan, &pekerjaan.Bidang_Industri,
        &pekerjaan.Lokasi_Kerja, &pekerjaan.Gaji_Range,
        &pekerjaan.Mulai_Kerja, &pekerjaan.Selesai_Kerja, &pekerjaan.Status_Pekerjaan,
		&pekerjaan.Jobdesk,
        &pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }

    _, err = db.Exec(`
        UPDATE pekerjaan_alumni
        SET is_deleted = FALSE, updated_at = NOW()
        WHERE id = $1 AND is_deleted = TRUE
    `, id)

    if err != nil {
        return nil, err
    }

    return &pekerjaan, nil
}

func PermanentDelete(db *sql.DB, id int) (*model.PekerjaanAlumni, error) {
	// Ambil data sebelum dihapus
	var pekerjaan model.PekerjaanAlumni
	row := db.QueryRow(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, 
		lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, 
		status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
		FROM pekerjaan_alumni WHERE id = $1
	`, id)

	err := row.Scan(
		&pekerjaan.ID, &pekerjaan.Alumni_ID, &pekerjaan.Nama_Perusahaan, 
		&pekerjaan.Posisi_Jabatan, &pekerjaan.Bidang_Industri,
		&pekerjaan.Lokasi_Kerja, &pekerjaan.Gaji_Range,
		&pekerjaan.Mulai_Kerja, &pekerjaan.Selesai_Kerja, 
		&pekerjaan.Status_Pekerjaan, &pekerjaan.Jobdesk,
		&pekerjaan.CreatedAt, &pekerjaan.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Query SQL untuk menghapus data
	_, err = db.Exec(`
		DELETE FROM pekerjaan_alumni WHERE id = $1 AND is_deleted = true
	`, id)

	if err != nil {
		return nil, err
	}

	return &pekerjaan, nil
}