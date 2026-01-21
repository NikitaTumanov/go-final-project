package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/NikitaTumanov/go-final-project/internal/pkg/db"
	"github.com/NikitaTumanov/go-final-project/internal/pkg/model"
)

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
