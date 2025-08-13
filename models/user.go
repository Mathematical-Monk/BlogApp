package models

type CreateUser struct {
	UserName     string `json:"username"`
	PasswordHash string `json:"-"`
}

type VerifyUser struct {
	// UserId   int64  `json:"userId"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
}
