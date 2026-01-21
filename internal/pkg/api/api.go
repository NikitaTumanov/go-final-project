package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NikitaTumanov/go-final-project/internal/pkg/db"
	"github.com/NikitaTumanov/go-final-project/internal/pkg/model"
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
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
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

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, errMethodNotAllowed.Error(), http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if nowStr == "" || dateStr == "" || repeat == "" {
		http.Error(w, errMissingQueryParameters.Error(), http.StatusBadRequest)
		return
	}

	now, err := time.Parse(dateFormat, nowStr)
	if err != nil {
		http.Error(w, errInvalidNowFormat.Error(), http.StatusBadRequest)
		return
	}

	next, err := nextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	decoder := json.NewDecoder(r.Body)

	var addTaskRequest model.AddTaskRequest
	err := decoder.Decode(&addTaskRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.AddTaskResponse{
			Error: errors.New("incorrect JSON").Error(),
		})
		return
	}

	if addTaskRequest.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.AddTaskResponse{
			Error: errors.New("title is required").Error(),
		})
		return
	}

	var date time.Time
	if strings.ReplaceAll(addTaskRequest.DateStr, " ", "") == "" {
		addTaskRequest.DateStr = time.Now().Format(dateFormat)
	} else {
		date, err = time.Parse(dateFormat, addTaskRequest.DateStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(model.AddTaskResponse{
				Error: errors.New("date is incorrect").Error(),
			})
			return
		}
		addTaskRequest.DateStr = date.Format(dateFormat)

		if date.Before(time.Now().Truncate(24 * time.Hour)) {
			if addTaskRequest.Repeate == "" {
				addTaskRequest.DateStr = time.Now().Format(dateFormat)
			} else {
				addTaskRequest.DateStr, err = nextDate(time.Now(), addTaskRequest.DateStr, addTaskRequest.Repeate)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					_ = json.NewEncoder(w).Encode(model.AddTaskResponse{
						Error: err.Error(),
					})
					return
				}
			}
		}
	}

	res, err := db.AddTask(addTaskRequest.DateStr, addTaskRequest.Title, addTaskRequest.Comment, addTaskRequest.Repeate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.AddTaskResponse{
			Error: errors.New("error in insert task to DB").Error(),
		})
		return
	}

	id, _ := res.LastInsertId()

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.AddTaskResponse{
		ID: strconv.FormatInt(id, 10),
	})
}

type tasksResponse struct {
	Tasks []model.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var tasks []model.Task
	search := r.URL.Query().Get("search")
	if search != "" {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			search = date.Format(dateFormat)
			tasks, err = db.TasksByDate(search, rowsLimit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(model.AddTaskResponse{
					Error: errors.New("error in select tasks from DB").Error(),
				})
				return
			}
		} else {
			tasks, err = db.TasksByTitleOrComment(search, rowsLimit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(model.AddTaskResponse{
					Error: errors.New("error in select tasks from DB").Error(),
				})
				return
			}
		}
	} else {
		var err error
		tasks, err = db.Tasks(rowsLimit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(model.AddTaskResponse{
				Error: errors.New("error in select tasks from DB").Error(),
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(tasksResponse{
		Tasks: tasks,
	})
}
