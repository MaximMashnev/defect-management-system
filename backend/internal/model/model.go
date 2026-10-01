package model

import (
	"time"
)

type User struct {
	ID           int       `json:"id"`
	Login        string    `json:"login"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type AuthUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

type Building struct {
	ID                 int       `json:"id"`
	Name               string    `json:"name"`
	Address            string    `json:"address"`
	DateRegistration   time.Time `json:"dateRegistration"`
	NumberApplications int       `json:"numberApplications"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
}

type Category struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type ApplicationStatus string

const (
	StatusNew        ApplicationStatus = "новый"
	StatusInProgress ApplicationStatus = "в работе"
	StatusCompleted  ApplicationStatus = "выполнено"
	StatusRejected   ApplicationStatus = "отклонено"
)

type Application struct {
	ID             int               `json:"id"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Email          string            `json:"email"`
	DateSubmission time.Time         `json:"dateSubmission"`
	DateInProgress *time.Time        `json:"dateInProgress,omitempty"`
	DateResult     *time.Time        `json:"dateResult,omitempty"`
	Status         ApplicationStatus `json:"status"`
	Priority       int               `json:"priority"`
	BuildingID     int               `json:"building_id"`
	UploadID       []int             `json:"upload_id"`
	UserID         *int              `json:"user_id,omitempty"`
	CategoryID     *int              `json:"category,omitempty"`
	Reason         *string           `json:"reason,omitempty"`
}

type UploadFile struct {
	ID           int       `json:"id"`
	OriginalName string    `json:"fileName"`
	StoredName   string    `json:"stored_name,omitempty"`
	FilePath     string    `json:"file_path,omitempty"`
	URL          string    `json:"url"`
	Bytes        int64     `json:"bytes"`
	Format       string    `json:"format"`
	MimeType     string    `json:"mime_type,omitempty"`
	Width        int       `json:"width"`
	Height       int       `json:"height"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
}

type Meta struct {
	CurrentPage    int `json:"current_page"`
	PerPage        int `json:"per_page"`
	RemainingCount int `json:"remaining_count"`
	TotalItems     int `json:"total_items"`
	TotalPages     int `json:"total_pages"`
}

type PaginatedApplications struct {
	Meta  Meta          `json:"meta"`
	Items []Application `json:"items"`
}

type PaginatedBuildings struct {
	Meta  Meta       `json:"meta"`
	Items []Building `json:"items"`
}
