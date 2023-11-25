package model

import "time"

type DbDevice struct {
	ID         int       `db:"id"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	MacAddress string    `db:"mac"`
	Name       string    `db:"name"`
	ActivePin  int       `db:"active_pin"`
}
