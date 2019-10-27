package models

type Item struct {
	Id     int    `db:"id"`
	NtinId int    `db:"ntinid"`
	Serial string `db:"serial"`
	Status int    `db:"status"`
}
