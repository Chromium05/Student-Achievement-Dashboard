package model

import "time"

type EmployedAlumni struct {
    Nama                 string    `db:"nama" json:"nama"`
    Jurusan              string    `db:"jurusan" json:"jurusan"`
    Angkatan             int       `db:"angkatan" json:"angkatan"`
    TahunLulus           int       `db:"tahun_lulus" json:"tahun_lulus"`
    NamaPerusahaan       string    `db:"nama_perusahaan" json:"nama_perusahaan"`
    LokasiKerja          string    `db:"lokasi_kerja" json:"lokasi_kerja"`
    BidangIndustri       string    `db:"bidang_industri" json:"bidang_industri"`
	PosisiJabatan 	     string    `db:"posisi_jabatan" json:"posisi_jabatan"`
    TanggalMulaiKerja  time.Time `db:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
    DeskripsiPekerjaan   string    `db:"deskripsi_pekerjaan" json:"deskripsi_pekerjaan"`

}
