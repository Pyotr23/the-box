-- +goose Up
-- +goose StatementBegin
CREATE TABLE countries (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL UNIQUE,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL
);

CREATE INDEX idx__countries__name ON countries(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx__countries__name;
DROP TABLE countries;
-- +goose StatementEnd
