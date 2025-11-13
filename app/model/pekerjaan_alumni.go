package model

import "time"

type PekerjaanAlumni struct {
	ID					int `json:"id"`
	Alumni_ID			int `json:"alumni_id"`
	Nama_Perusahaan		string `json:"nama_perusahaan"`
	Posisi_Jabatan		string `json:"posisi_jabatan"`
	Bidang_Industri		string `json:"bidang_industri"`
	Lokasi_Kerja 		string `json:"lokasi_kerja"`
	Gaji_Range 			string `json:"gaji_range"`
	Mulai_Kerja			*time.Time `json:"tanggal_mulai_kerja"`
	Selesai_Kerja 		*time.Time `json:"tanggal_selesai_kerja"`
	Status_Pekerjaan	string `json:"status_pekerjaan"` // Contoh: bekerja, selesai, resigned
	Jobdesk				string `json:"deskripsi_pekerjaan"`
	CreatedAt			time.Time `json:"created_at"`
	UpdatedAt			time.Time `json:"updated_at"`
	IsDeleted			bool `json:"is_deleted"`
}

type CreatePekerjaan struct {
	Alumni_ID			int `json:"alumni_id"`
	Nama_Perusahaan		string `json:"nama_perusahaan"`
	Posisi_Jabatan		string `json:"posisi_jabatan"`
	Bidang_Industri		string `json:"bidang_industri"`
	Lokasi_Kerja 		string `json:"lokasi_kerja"`
	Gaji_Range 			string `json:"gaji_range"`
	Mulai_Kerja			string `json:"tanggal_mulai_kerja"`
	Selesai_Kerja 		string `json:"tanggal_selesai_kerja"`
	Status_Pekerjaan	string `json:"status_pekerjaan"` // Contoh: bekerja, selesai, resigned
	Jobdesk				string `json:"deskripsi_pekerjaan"`
	UpdatedAt			time.Time `json:"updated_at"`
}

type SoftDeletePekerjaanAlumni struct {
	ID					int `json:"id"`
	Alumni_ID			int `json:"alumni_id"`
	Nama_Perusahaan		string `json:"nama_perusahaan"`
	Posisi_Jabatan		string `json:"posisi_jabatan"`
	Bidang_Industri		string `json:"bidang_industri"`
	Lokasi_Kerja 		string `json:"lokasi_kerja"`
	Gaji_Range 			string `json:"gaji_range"`
	Mulai_Kerja			*time.Time `json:"tanggal_mulai_kerja"`
	Selesai_Kerja 		*time.Time `json:"tanggal_selesai_kerja"`
	Status_Pekerjaan	string `json:"status_pekerjaan"` // Contoh: bekerja, selesai, resigned
	Jobdesk				string `json:"deskripsi_pekerjaan"`
	CreatedAt			time.Time `json:"created_at"`
	UpdatedAt			time.Time `json:"updated_at"`
	IsDeleted			bool `json:"id_deleted"`
}