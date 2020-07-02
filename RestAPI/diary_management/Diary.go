package diary_management

import (
	uuid2 "github.com/satori/go.uuid"
	"time"
)

type Diary struct {
	Id        uuid2.UUID    `json:"id"`
	Title     string    	`json:"title"`
	Body      string    	`json:"content"`
	UserId    uuid2.UUID	`json:"user_id"`
	CreatedAt time.Time 	`json:"created_at"`
	UpdatedAt time.Time 	`json:"updated_at"`
}

type UpdateDiary struct {
	Id        uuid2.UUID    `json:"id"`
	Title     string    	`json:"title"`
	Body      string    	`json:"content"`
	UpdatedAt time.Time 	`json:"updated_at"`
}

type GetQuarterlyDiary struct{
	Year		int		`json:"year"`
	Quarter		int 	`json:"Quarter"`
}

const (
	layoutISO = "2006-01-02 00:00:00"
)