package model

type Alumni struct {
    Nama       *string `json:"nama"`
    IDFakultas *int    `json:"id_fakultas"`
    IDProdi    *int    `json:"id_prodi"`
    TahunLulus *int    `json:"tahun_lulus"`
    Sumber     *string `json:"sumber"`
    IDSumber   *string `json:"id_sumber"`
}

type SemuaAlumni struct {
    NIM        *string `json:"nim"`
    Nama       *string `json:"nama"`
    IDFakultas *int    `json:"id_fakultas"`
    IDProdi    *int    `json:"id_prodi"`
    TahunLulus *int    `json:"tahun_lulus"`
    Sumber     *string `json:"sumber"`
    IDSumber   *string `json:"id_sumber"`
}