package models

import "time"

type LockBoxInput struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Login       string `json:"login"`
	Password    string `json:"password"`
	Description string `json:"description"`
}

type LockBox struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Login       string    `json:"login"`
	Password    string    `json:"password"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SyncedAt    time.Time `json:"synced_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}

type DeleteLockBox struct {
	ID     int
	ItemId int    `json:"item_id"`
	Name   string `json:"name"`
}
type UpdatedItem struct {
	ID   int
	Name string
}
