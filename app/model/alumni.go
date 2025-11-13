package model

import "time"

type Alumni struct {
    ID          int    `json:"id"`
    NIM         string `json:"nim"`
    Nama        string `json:"nama"`
    Jurusan     string `json:"jurusan"`
    Angkatan    int    `json:"angkatan"`
    Tahun_Lulus int    `json:"tahun_lulus"`
    Email       string `json:"email"`
    No_Telepon  string `json:"no_telepon"`
    Alamat      string `json:"alamat"`
    CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAlumni struct {
    ID          int    
    NIM         string `json:"nim"`
    Nama        string `json:"nama"`
    Jurusan     string `json:"jurusan"`
    Angkatan    int    `json:"angkatan"`
    Tahun_Lulus int    `json:"tahun_lulus"`
    Email       string `json:"email"`
    No_Telepon  string `json:"no_telepon"`
    Alamat      string `json:"alamat"`
}

type SoftDeleteAlumni struct {
    ID          int    `json:"id"`
    NIM         string `json:"nim"`
    Nama        string `json:"nama"`
    Jurusan     string `json:"jurusan"`
    Angkatan    int    `json:"angkatan"`
    Tahun_Lulus int    `json:"tahun_lulus"`
    Email       string `json:"email"`
    No_Telepon  string `json:"no_telepon"`
    Alamat      string `json:"alamat"`
    CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
    IsDeleted   bool    `json:"is_deleted"`
}