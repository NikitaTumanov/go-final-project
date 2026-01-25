package api

import (
	"encoding/json"
	"errors"
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
		_ = json.NewEncoder(w).Encode(model.GetTaskByIDResponse{
			Error: errors.New("ID is empty").Error(),
		})
		return
	}

	if _, err := strconv.Atoi(idStr); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(model.GetTaskByIDResponse{
			Error: errors.New("ID is incorrect").Error(),
		})
		return
	}

	task, err := db.TasksByID(idStr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(model.GetTaskByIDResponse{
			Error: errors.New("error in select task from DB").Error(),
		})
		return
	}

	if task.Repeate == "" {
		_, err = db.DeleteTaskByID(idStr)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(model.UpdateTaskByIDResponse{
				Error: errors.New("error in delete task from DB").Error(),
			})
			return
		}
	} else {
		task.Date, err = nextDate(time.Now(), task.Date, task.Repeate)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(model.UpdateTaskByIDResponse{
				Error: err.Error(),
			})
			return
		}

		_, err = db.ChangeTaskByID(task.ID, task.Date, task.Title, task.Comment, task.Repeate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(model.UpdateTaskByIDResponse{
				Error: errors.New("error in update task to DB").Error(),
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.UpdateTaskByIDResponse{})
}
