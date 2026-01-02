package workout

import (
	"github.com/jmoiron/sqlx"
)

type Repository interface {
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}
