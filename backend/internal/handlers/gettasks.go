package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ilexsor/internal/models"
	"gorm.io/gorm"
)

// Обработчик для GET /api/tasks
func GetTasks(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var task []models.Task
		response := models.TasksResponse{
			Tasks: []models.Task{},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		tx := db.WithContext(ctx).Limit(models.TaskLimit).Find(&task)

		if tx.Error != nil {

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.GetTasksError,
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)
			return
		}

		response.Tasks = task
		resp, err := json.Marshal(response)

		if err != nil {

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.MarshalError,
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}
