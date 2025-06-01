package models

import "time"

type Data struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	Url         string     `json:"url"`
	Login       string     `json:"login"`
	Password    string     `json:"password"`
	Description string     `json:"description"`
	UserID      int        `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
