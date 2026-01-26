package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/NikitaTumanov/go-final-project/internal/pkg/db"
	"github.com/NikitaTumanov/go-final-project/internal/pkg/model"
)

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tasks := make([]model.Task, 0)
	search := r.URL.Query().Get("search")
	if search != "" {
		if date, err := time.Parse("02.01.2006", search); err == nil {
			search = date.Format(dateFormat)
			tasks, err = db.TasksByDate(search, rowsLimit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(model.GetTasksResponse{
					Error: errDatabaseSelect.Error(),
				})
				return
			}
		} else {
			tasks, err = db.TasksByTitleOrComment(search, rowsLimit)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(model.GetTasksResponse{
					Error: errDatabaseSelect.Error(),
				})
				return
			}
		}
	} else {
		var err error
		tasks, err = db.Tasks(rowsLimit)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(model.GetTasksResponse{
				Error: errDatabaseSelect.Error(),
			})
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(model.GetTasksResponse{
		Tasks: tasks,
	})
}
