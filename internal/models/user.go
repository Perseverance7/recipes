package models

type User struct {
	Id       int    `db:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Salt     string `db:"salt"`
}
