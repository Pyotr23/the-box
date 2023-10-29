package repo

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	driverName = "sqlite3"
	dbSource   = "./bluetooth-api/db/bluetooth.db"
)

type Repository struct {
}

func NewRepo() (Repository, error) {
	db, err := sql.Open(driverName, dbSource)
	if err != nil {
		return Repository{}, fmt.Errorf("open: %w", err)
	}

	if err = db.Ping(); err != nil {
		return Repository{}, fmt.Errorf("ping: %w", err)
	}

	return Repository{}, nil
}
