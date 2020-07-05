package diary_management

type Diary struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"content"`
	UserId    int64     `json:"user_id"`
	CreatedAt string 	`json:"created_at"`
	UpdatedAt string 	`json:"updated_at"`
}

type UpdateDiary struct {
	Id        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"content"`
	UpdatedAt string `json:"updated_at"`
}

type GetQuarterlyDiary struct {
	Year    int `json:"year"`
	Quarter int `json:"Quarter"`
}

const (
	layoutISO = "2006-01-02 00:00:00"
)
