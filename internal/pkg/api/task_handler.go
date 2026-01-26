package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/NikitaTumanov/go-final-project/internal/pkg/db"
	"github.com/NikitaTumanov/go-final-project/internal/pkg/model"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	decoder := json.NewDecoder(r.Body)

	var addTaskRequest model.AddTaskRequest
	err := decoder.Decode(&addTaskRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errIncorrectJSON.Error(),
		})
		return
	}

	if addTaskRequest.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errEmptyTitle.Error(),
		})
		return
	}

	var date time.Time
	if addTaskRequest.DateStr == "" {
		addTaskRequest.DateStr = time.Now().Format(dateFormat)
	} else {
		date, err = time.Parse(dateFormat, addTaskRequest.DateStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(model.ErrorResponse{
				Error: errIncorrectDate.Error(),
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
					_ = json.NewEncoder(w).Encode(model.ErrorResponse{
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
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errDatabaseInsert.Error(),
		})
		return
	}

	id, _ := res.LastInsertId()

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.AddTaskResponse{
		ID: strconv.FormatInt(id, 10),
	})
}

func getTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errEmptyId.Error(),
		})
		return
	}

	if _, err := strconv.Atoi(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errIncorrectId.Error(),
		})
		return
	}

	task, err := db.TaskByID(idStr)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(model.ErrorResponse{
				Error: errDatabaseNoTasks.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errDatabaseSelect.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.GetTaskByIDResponse{
		ID:      task.ID,
		Date:    task.Date,
		Title:   task.Title,
		Comment: task.Comment,
		Repeate: task.Repeate,
	})
}

func changeTaskByIDHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	decoder := json.NewDecoder(r.Body)

	var task model.Task
	err := decoder.Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errIncorrectJSON.Error(),
		})
		return
	}

	if task.ID == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errEmptyId.Error(),
		})
		return
	}

	if _, err := strconv.Atoi(task.ID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errEmptyId.Error(),
		})
		return
	}

	if task.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errEmptyTitle.Error(),
		})
		return
	}

	var date time.Time
	if task.Date == "" {
		task.Date = time.Now().Format(dateFormat)
	} else {
		date, err = time.Parse(dateFormat, task.Date)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(model.ErrorResponse{
				Error: errIncorrectDate.Error(),
			})
			return
		}
		task.Date = date.Format(dateFormat)

		if date.Before(time.Now().Truncate(24 * time.Hour)) {
			if task.Repeate == "" {
				task.Date = time.Now().Format(dateFormat)
			} else {
				task.Date, err = nextDate(time.Now(), task.Date, task.Repeate)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					_ = json.NewEncoder(w).Encode(model.ErrorResponse{
						Error: err.Error(),
					})
					return
				}
			}
		}
	}

	res, err := db.ChangeTaskByID(task.ID, task.Date, task.Title, task.Comment, task.Repeate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errDatabaseUpdate.Error(),
		})
		return
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errDatabaseUpdate.Error(),
		})
		return
	}

	if affectedRows == 0 {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errDatabaseNoTasks.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.ErrorResponse{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errEmptyId.Error(),
		})
		return
	}

	if _, err := strconv.Atoi(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errIncorrectId.Error(),
		})
		return
	}

	_, err := db.DeleteTaskByID(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.ErrorResponse{
			Error: errDatabaseDelete.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.ErrorResponse{})
}
