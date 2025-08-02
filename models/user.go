package models

type CreateUser struct {
	Usename      string `json:"username"`
	PasswordHash string `json:"-"`
}


type User struct {
	Id uint64 `json:"id"`
	Username string `json:"username"`
}