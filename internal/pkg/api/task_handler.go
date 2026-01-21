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
