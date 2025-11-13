package repository

import (
	"crud-alumni-5/app/model"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Custom error untuk duplicate key
var ErrDuplicateNIM = errors.New("NIM sudah terdaftar")
var ErrDuplicateEmail = errors.New("email sudah terdaftar")
var ErrInvalidInput = errors.New("input tidak valid")

// GET semua alumni (GET 127.0.0.1:3000/:key)
func GetAllAlumni(db *sql.DB, search, sortBy, order string, limit, offset int) ([]model.Alumni, error) {
	query := fmt.Sprintf(`
		SELECT id, nim, nama, jurusan, angkatan, tahun_lulus, email, 
		no_telepon, alamat, created_at, updated_at
		FROM alumni
		WHERE nama ILIKE $1 OR nim ILIKE $1 OR jurusan ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
		`, sortBy, order)

	rows, err := db.Query(query, "%"+search+"%", limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var alumniList []model.Alumni

	for rows.Next() {
		var alumni model.Alumni
		err := rows.Scan(
			&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan,
			&alumni.Angkatan, &alumni.Tahun_Lulus,
			&alumni.Email, &alumni.No_Telepon, &alumni.Alamat,
			&alumni.CreatedAt, &alumni.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		alumniList = append(alumniList, alumni)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return alumniList, nil
}

// GET alumni by ID (GET 127.0.0.1:3000/:key/:id)
func GetAlumniByID(db *sql.DB, id int) (*model.Alumni, error) {
	// Generate query Postgresql
	row := db.QueryRow(`
		SELECT id, nim, nama, jurusan, angkatan, 
		tahun_lulus, email, no_telepon, alamat, 
		created_at, updated_at
		FROM alumni
		WHERE id = $1
	`, id)

	var alumni model.Alumni
	err := row.Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan,
		&alumni.Tahun_Lulus, &alumni.Email, &alumni.No_Telepon, &alumni.Alamat,
		&alumni.CreatedAt, &alumni.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Alumni tidak ada
		}
		return nil, err // Error lainnya
	}

	return &alumni, nil
}

// Fungsi helper untuk cek apakah NIM sudah ada
func checkNIMExists(db *sql.DB, nim string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM alumni WHERE nim = $1", nim).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Fungsi helper untuk cek apakah Email sudah ada
func checkEmailExists(db *sql.DB, email string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM alumni WHERE email = $1", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// POST tambah data alumni baru (POST 127.0.0.1:3000/:key)
func PostNewAlumni(db *sql.DB, alumni *model.CreateAlumni) (*model.Alumni, error) {
	// Validasi input kosong
	if strings.TrimSpace(alumni.NIM) == "" || strings.TrimSpace(alumni.Nama) == "" ||
		strings.TrimSpace(alumni.Jurusan) == "" || alumni.Angkatan == 0 ||
		alumni.Tahun_Lulus == 0 || strings.TrimSpace(alumni.Email) == "" ||
		strings.TrimSpace(alumni.No_Telepon) == "" || strings.TrimSpace(alumni.Alamat) == "" {
		return nil, ErrInvalidInput
	}

	// Cek apakah NIM sudah ada
	nimExists, err := checkNIMExists(db, alumni.NIM)
	if err != nil {
		return nil, err
	}
	if nimExists {
		return nil, ErrDuplicateNIM
	}

	// Cek apakah Email sudah ada (opsional, tergantung apakah email juga harus unique)
	emailExists, err := checkEmailExists(db, alumni.Email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrDuplicateEmail
	}

	// Insert data baru
	var newID int
	err = db.QueryRow(`
		INSERT INTO alumni (nim, nama, jurusan, angkatan, 
		tahun_lulus, email, no_telepon, alamat, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id
	`, strings.TrimSpace(alumni.NIM), strings.TrimSpace(alumni.Nama),
		strings.TrimSpace(alumni.Jurusan), alumni.Angkatan, alumni.Tahun_Lulus, strings.TrimSpace(alumni.Email),
		strings.TrimSpace(alumni.No_Telepon), strings.TrimSpace(alumni.Alamat)).Scan(&newID)

	if err != nil {
		// Meng-handle input NIM dan Email yang duplikat
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			if strings.Contains(err.Error(), "alumni_pkey") || strings.Contains(err.Error(), "nim") {
				return nil, ErrDuplicateNIM
			}
			if strings.Contains(err.Error(), "email") {
				return nil, ErrDuplicateEmail
			}
		}
		return nil, err
	}

	// Ambil data alumni yang baru ditambahkan
	var newAlumni model.Alumni
	row := db.QueryRow(`
		SELECT id, nim, nama, jurusan, angkatan, 
		tahun_lulus, email, no_telepon, alamat, 
		created_at, updated_at
		FROM alumni WHERE id = $1
	`, newID)

	err = row.Scan(
		&newAlumni.ID, &newAlumni.NIM, &newAlumni.Nama, &newAlumni.Jurusan, &newAlumni.Angkatan,
		&newAlumni.Tahun_Lulus, &newAlumni.Email, &newAlumni.No_Telepon, &newAlumni.Alamat,
		&newAlumni.CreatedAt, &newAlumni.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &newAlumni, nil
}

// PUT Update data alumni yang sudah ada (PUT 127.0.0.1:3000/:key)
func UpdateAlumni(db *sql.DB, alumni *model.CreateAlumni, id int) (*model.Alumni, error) {
	// Validasi input kosong
	if strings.TrimSpace(alumni.NIM) == "" || strings.TrimSpace(alumni.Nama) == "" ||
		strings.TrimSpace(alumni.Jurusan) == "" || alumni.Angkatan == 0 ||
		alumni.Tahun_Lulus == 0 || strings.TrimSpace(alumni.Email) == "" ||
		strings.TrimSpace(alumni.No_Telepon) == "" || strings.TrimSpace(alumni.Alamat) == "" {
		return nil, ErrInvalidInput
	}

	_, err := db.Exec(`
		UPDATE alumni 
		SET nama = $1, nim = $2, jurusan = $3, angkatan = $4, 
		tahun_lulus = $5, email = $6, no_telepon = $7,
		alamat = $8, updated_at = $9 
		WHERE id = $10`, alumni.Nama, alumni.NIM, alumni.Jurusan, alumni.Angkatan,
		alumni.Tahun_Lulus, alumni.Email,
		alumni.No_Telepon, alumni.Alamat, time.Now(), id)

	if err != nil {
		return nil, err
	}

	// Ambil data alumni yang baru diupdate
	var updatedAlumni model.Alumni
	row := db.QueryRow(`
		SELECT id, nim, nama, jurusan, angkatan, 
		tahun_lulus, email, no_telepon, alamat, 
		created_at, updated_at
		FROM alumni WHERE id = $1
	`, id)
	errScan := row.Scan(
		&updatedAlumni.ID, &updatedAlumni.NIM, &updatedAlumni.Nama,
		&updatedAlumni.Jurusan, &updatedAlumni.Angkatan,
		&updatedAlumni.Tahun_Lulus, &updatedAlumni.Email,
		&updatedAlumni.No_Telepon, &updatedAlumni.Alamat,
		&updatedAlumni.CreatedAt, &updatedAlumni.UpdatedAt,
	)

	if errScan != nil {
		return nil, errScan
	}

	return &updatedAlumni, nil
}

func DeleteAlumni(db *sql.DB, id int) (*model.Alumni, error) {
	// Ambil data sebelum dihapus
	var alumni model.Alumni
	row := db.QueryRow(`
		SELECT id, nim, nama, jurusan, angkatan, 
		tahun_lulus, email, no_telepon, alamat, created_at, updated_at
		FROM alumni WHERE id = $1
	`, id)
	errScan := row.Scan(
		&alumni.ID, &alumni.NIM, &alumni.Nama, &alumni.Jurusan, &alumni.Angkatan,
		&alumni.Tahun_Lulus, &alumni.Email, &alumni.No_Telepon, &alumni.Alamat,
		&alumni.CreatedAt, &alumni.UpdatedAt,
	)

	if errScan != nil {
		return nil, errScan
	}

	// Query SQL untuk menghapus data
	_, err := db.Exec(`
		DELETE FROM alumni WHERE id = $1
	`, id)

	if err != nil {
		return nil, err
	}

	return &alumni, nil
}
