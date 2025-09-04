package repository

import (
	"database/sql"
	"tugasminggu3/app/model"
)

func CheckAlumniByNim(db *sql.DB, nim string) (*model.Alumni, error) {
	alumni := new(model.Alumni)
	query := `SELECT nama, 1 as idfakultas, 1 as idprodi, 2020 as tahun_lulus,
              'Alumnipedia' as sumber, '123' as id_sumber
              FROM alumnipedia WHERE nim = $1 LIMIT 1`

	err := db.QueryRow(query, nim).Scan(
		&alumni.Nama,
		&alumni.IDFakultas,
		&alumni.IDProdi,
		&alumni.TahunLulus,
		&alumni.Sumber,
		&alumni.IDSumber,
	)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

func GetAllAlumni(db *sql.DB) (*model.SemuaAlumni, error) {
	model := new(model.SemuaAlumni)
	query := `SELECT * FROM alumnipedia`

	err := db.QueryRow(query).Scan(
		&model.NIM,
		&model.Nama,
		&model.IDFakultas,
		&model.IDProdi,
		&model.TahunLulus,
		&model.Sumber,
		&model.IDSumber,
	)

	if err != nil {
		return nil, err
	}
	return model, nil
}
