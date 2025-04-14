package models

type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Done      bool   `json:"done"`
	CreatedAt string `json:"createdAt"`
	DueDate   string `json:"dueDate"` // örn: "2025-04-30"
}
