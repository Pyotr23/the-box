package model

type Device struct {
	ID         int    `db:"id"`
	MacAddress string `db:"mac"`
	Name       string `db:"name"`
}
