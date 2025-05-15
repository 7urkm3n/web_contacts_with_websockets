package models

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Contacts ContactModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Contacts: ContactModel{DB: db},
	}
}
