package api

import (
	"errors"
	"net/http"
)

const (
	dateFormat = "20060102"
	rowsLimit  = 50
)

var (
	errMethodNotAllowed       = errors.New("method not allowed")
	errMissingQueryParameters = errors.New("missing query parameters")
	errInvalidNowFormat       = errors.New("invalid now format")
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", taskHandler)
	http.HandleFunc("/api/tasks", tasksHandler)
	http.HandleFunc("/api/task/done", doneTaskHandler)
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskByIDHandler(w, r)
	case http.MethodPut:
		changeTaskByIDHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTasksHandler(w, r)
	default:
		http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
	}
}
