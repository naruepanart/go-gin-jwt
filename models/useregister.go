package models

type UserRegister struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}