package models

type Record struct {
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}
