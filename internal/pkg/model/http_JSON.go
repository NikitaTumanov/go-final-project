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
	ID string `json:"id"`
}

type GetTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type GetTaskByIDResponse struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeate string `json:"repeat"`
}

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}
