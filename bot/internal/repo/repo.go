package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Pyotr23/the-box/bot/internal/pkg/model"
	_ "github.com/mattn/go-sqlite3"
)

const (
	driverName = "sqlite3"
	dbSource   = "./bot/db/bot.db"
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

func InsertJob(ctx context.Context, data model.JobSettingsChatID) error {
	bytesSettings, err := json.Marshal(data.JobSettings)
	if err != nil {
		return fmt.Errorf("marshal job settings: %w", err)
	}

	const q = `
		insert into device(chat_id, settings_json)
		values (?, ?)`

	return exec(ctx, q, data.ChatID, string(bytesSettings))
}

func GetJobs(ctx context.Context) ([]model.JobSettingsChatID, error) {
	const q = `
		select chat_id
		  , settings_json
		from job`

	rows, err := query(ctx, q)
	defer func() {
		if dErr := rows.Close(); dErr != nil {
			log.Printf("rows close: %s", dErr)
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	var res []model.JobSettingsChatID
	for rows.Next() {
		var (
			chatID  int64
			strItem string
		)
		if err = rows.Scan(chatID, &chatID, &strItem); err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		var js model.JobSettings
		if err = json.Unmarshal([]byte(strItem), &js); err != nil {
			return nil, fmt.Errorf("unmarshal job settings")
		}

		res = append(res, model.JobSettingsChatID{
			ChatID:      chatID,
			JobSettings: js,
		})
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
