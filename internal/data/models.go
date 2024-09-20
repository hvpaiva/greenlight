package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

const (
	PermissionMovieRead  = "movies:read"
	PermissionMovieWrite = "movies:write"
)

type Models struct {
	Movies     MovieModel
	Users      UserModel
	Tokens     TokenModel
	Permission PermissionModel
}

func New(db *sql.DB) *Models {
	return &Models{
		Movies:     MovieModel{DB: db},
		Users:      UserModel{DB: db},
		Tokens:     TokenModel{DB: db},
		Permission: PermissionModel{DB: db},
	}
}
