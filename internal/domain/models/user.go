package models

type User struct {
	Uuid     string `json:"uuid"`
	Email    string `json:"email"`
	PassHash []byte `json:"password"`
}
