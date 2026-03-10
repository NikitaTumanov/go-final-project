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
	errIncorrectJSON          = errors.New("incorrect JSON")

	errEmptyId            = errors.New("id is empty")
	errIncorrectId        = errors.New("id is incorrect")
	errEmptyTitle         = errors.New("title is empty")
	errIncorrectNowFormat = errors.New("invalid now format")
	errIncorrectRepeat    = errors.New("incorrect repeat format")
	errIncorrectDate      = errors.New("date is incorrect")
	errParseStartDate     = errors.New("start date parse error")

	errNextDateNotFound = errors.New("next date not found")

	errDatabaseSelect  = errors.New("error in select task from DB")
	errDatabaseInsert  = errors.New("error in insert task to DB")
	errDatabaseUpdate  = errors.New("error in update task to DB")
	errDatabaseDelete  = errors.New("error in delete task from DB")
	errDatabaseNoTasks = errors.New("no tasks with this id")

	errWrongPasswordJSON = errors.New("wrong password JSON")
	errWrongPassword     = errors.New("wrong password")
	errCreateJWT         = errors.New("can not create JWT")
)

func Init() {
	http.HandleFunc("/api/nextdate", nextDayHandler)
	http.HandleFunc("/api/task", auth(taskHandler))
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(doneTaskHandler))
	http.HandleFunc("/api/signin", signInHandler)
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
