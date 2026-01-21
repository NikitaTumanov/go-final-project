package model

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeate string `json:"repeat"`
}

type AddTaskRequest struct {
	DateStr string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeate string `json:"repeat"`
}

type AddTaskResponse struct {
	ID    string `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}
