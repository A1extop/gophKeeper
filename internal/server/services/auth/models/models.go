package models

type AuthUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type InfoUser struct {
	UserId   int
	Username string
	UserType string
}
