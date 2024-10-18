package models

type User struct {
	Id       int64  `json:"id"`
	Email    string `json:"email"`
	PassHash []byte `json:"password"`
}
