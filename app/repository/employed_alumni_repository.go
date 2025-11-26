package repository

import (
	"student-report/app/model"
	"database/sql"
)

func GetAllEmployedAlumni(db *sql.DB) ([]model.EmployedAlumni, error) {
	rows, err := db.Query(`
		SELECT nama, jurusan, angkatan, tahun_lulus, nama_perusahaan, lokasi_kerja, bidang_industri, posisi_jabatan,
		tanggal_mulai_kerja, deskripsi_pekerjaan
		FROM employed_alumni
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var employedAlumniList []model.EmployedAlumni

	for rows.Next() {
		var employed model.EmployedAlumni
		err := rows.Scan(
			&employed.Nama, &employed.Jurusan, &employed.Angkatan, &employed.TahunLulus,
			&employed.NamaPerusahaan, &employed.LokasiKerja, &employed.BidangIndustri, &employed.PosisiJabatan,
			&employed.TanggalMulaiKerja, &employed.DeskripsiPekerjaan,
		)
		if err != nil {
			return nil, err
		}
		employedAlumniList = append(employedAlumniList, employed)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employedAlumniList, nil
}

func GetEmployedAlumniLessThreeYears(db *sql.DB) ([]model.EmployedAlumni, error) {
	rows, err := db.Query(`
		SELECT nama, jurusan, angkatan, tahun_lulus, nama_perusahaan, bidang_industri, posisi_jabatan,
		tanggal_mulai_kerja, deskripsi_pekerjaan
		FROM employed_alumni
		WHERE tanggal_mulai_kerja > NOW() - INTERVAL '3 YEARS'
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var employedAlumniList []model.EmployedAlumni

	for rows.Next() {
		var employed model.EmployedAlumni
		err := rows.Scan(
			&employed.Nama, &employed.Jurusan, &employed.Angkatan, &employed.TahunLulus,
			&employed.NamaPerusahaan, &employed.BidangIndustri, &employed.PosisiJabatan,
			&employed.TanggalMulaiKerja, &employed.DeskripsiPekerjaan,
		)
		if err != nil {
			return nil, err
		}
		employedAlumniList = append(employedAlumniList, employed)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return employedAlumniList, nil
}
