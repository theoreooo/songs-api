package models

import (
	"time"
)

type Song struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GroupName   string    `gorm:"column:group_name" json:"group"`
	Song        string    `json:"song"`
	ReleaseDate string    `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
