package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/Pyotr23/the-box/bluetooth-api/internal/pkg/model"
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

func (Repository) UpsertDevice(ctx context.Context, name, macAddress string) error {
	const q = `
		insert into device (mac, name)
		values (?, ?)
		on conflict do update 
		set name = ?
			, updated_at = current_timestamp`

	return exec(ctx, q, macAddress, name, name)
}

func (Repository) DeleteDevice(ctx context.Context, id int) error {
	const q = `
		delete from device
		where id = ?`

	return exec(ctx, q, id)
}

func (Repository) GetByMacAddresses(ctx context.Context, macAddresses []string) ([]model.DbDevice, error) {
	if len(macAddresses) == 0 {
		return nil, nil
	}

	placeholderString, genericMacAddresses := getQueryInfo[string](macAddresses)
	q := fmt.Sprintf(`
			select *
			from device
			where mac in (%s)`,
		placeholderString,
	)

	rows, err := query(ctx, q, genericMacAddresses...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer func() {
		if dErr := rows.Close(); dErr != nil {
			log.Printf("rows close: %s", dErr)
		}
	}()

	var res = make([]model.DbDevice, 0, len(macAddresses))
	for rows.Next() {
		var device = model.DbDevice{}
		rErr := rows.Scan(
			&device.ID,
			&device.CreatedAt,
			&device.UpdatedAt,
			&device.MacAddress,
			&device.Name,
		)
		if rErr != nil {
			log.Fatal("rows scan error")
		}

		res = append(res, device)
	}

	return res, nil
}

func getDB() (*sql.DB, error) {
	return sql.Open(driverName, dbSource)
}

func tryCloseDB(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("close connection: %s", err)
	}
}

func exec(ctx context.Context, query string, args ...any) error {
	stm, err := getPreparedStatement(ctx, query)
	if err != nil {
		return fmt.Errorf("get prepared statement: %w", err)
	}

	defer func() {
		if dErr := stm.Close(); dErr != nil {
			log.Printf("close prepared statement in exec: %s", err)
		}
	}()

	if _, err = stm.ExecContext(ctx, args...); err != nil {
		return fmt.Errorf("exec context: %w", err)
	}

	return nil
}

func query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	stm, err := getPreparedStatement(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get prepared statement: %w", err)
	}

	defer func() {
		if dErr := stm.Close(); dErr != nil {
			log.Printf("close prepared statement in exec: %s", err)
		}
	}()

	rows, err := stm.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}

	return rows, nil
}

func getPreparedStatement(ctx context.Context, query string) (*sql.Stmt, error) {
	db, err := getDB()
	if err != nil {
		return nil, fmt.Errorf("get db: %w", err)
	}

	stm, err := db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("prepare context: %w", err)
	}

	return stm, nil
}

func getQueryInfo[T any](sl []T) (string, []any) {
	if len(sl) == 0 {
		return "", nil
	}

	var (
		anySl        = make([]any, 0, len(sl))
		placeholders = make([]string, len(sl))
	)
	for i, item := range sl {
		placeholders[i] = "?"
		anySl = append(anySl, item)
	}

	return strings.Join(placeholders, ","), anySl
}
