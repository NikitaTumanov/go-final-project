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

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if task.Repeate == "" {
		_, err = db.DeleteTaskByID(idStr)
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
				Error: errDatabaseDelete.Error(),
			})
			return
		}
	} else {
		task.Date, err = nextDate(time.Now(), task.Date, task.Repeate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(model.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		_, err = db.ChangeTaskByID(task.ID, task.Date, task.Title, task.Comment, task.Repeate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(model.ErrorResponse{
				Error: errDatabaseUpdate.Error(),
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.ErrorResponse{})
}
