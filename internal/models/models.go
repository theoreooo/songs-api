package models

import (
	"time"
)

type Song struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ArtistID    uint      `gorm:"index" json:"artistId"`
	Artist      Artist    `gorm:"foreignKey:ArtistID" json:"artist"`
	Song        string    `gorm:"not null" json:"song"`
	ReleaseDate time.Time `json:"releaseDate" example:"2025-01-16"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Artist struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null" json:"group"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SongUpdate struct {
	GroupName   *string    `json:"group,omitempty"`
	Song        *string    `json:"song,omitempty"`
	ReleaseDate *time.Time `json:"releaseDate,omitempty" example:"2025-01-16"`
	Text        *string    `json:"text,omitempty"`
	Link        *string    `json:"link,omitempty"`
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
