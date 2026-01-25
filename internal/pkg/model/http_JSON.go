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

type GetTasksResponse struct {
	Tasks []Task `json:"tasks"`
	Error string `json:"error,omitempty"`
}

type GetTaskByIDResponse struct {
	ID      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeate string `json:"repeat,omitempty"`
	Error   string `json:"error,omitempty"`
}

type UpdateTaskByIDResponse struct {
	Error string `json:"error,omitempty"`
}
