package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	driverName = "sqlite3"
	dbSource   = "./bluetooth-api/db/bluetooth.db"
)

type Repository struct{}

func NewRepo() (Repository, error) {
	db, err := getDB()
	if err != nil {
		return Repository{}, fmt.Errorf("get db: %w", err)
	}

	defer tryCloseDB(db)

	if err = db.Ping(); err != nil {
		return Repository{}, fmt.Errorf("ping: %w", err)
	}

	db.Close()

	return Repository{}, nil
}

func (_ Repository) UpsertDevice(ctx context.Context, name, macAddress string) error {
	const q = `
		insert into device (mac, name)
		values (?, ?)
		on conflict do update 
		set name = ?
			, updated_at = ?`

	db, err := getDB()
	if err != nil {
		return fmt.Errorf("get db: %w", err)
	}

	stm, err := db.PrepareContext(ctx, q)
	if err != nil {
		return fmt.Errorf("prepare context: %w", err)
	}

	if _, err = stm.ExecContext(ctx, macAddress, name, name, time.Now()); err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	return nil
}

func getDB() (*sql.DB, error) {
	return sql.Open(driverName, dbSource)
}

func tryCloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("close connection: %s", err)
	}
}
