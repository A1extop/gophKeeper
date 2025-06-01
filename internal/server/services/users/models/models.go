package models

import "time"

type User struct {
	UserId    int       `json:"user_id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	UserType  string    `json:"user_type"` // тип user_type ('admin', attendee');
	CreatedAt time.Time `json:"-"`
}
