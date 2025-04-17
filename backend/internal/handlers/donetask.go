package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ilexsor/internal/models"
	"github.com/ilexsor/internal/utils"
	"gorm.io/gorm"
)

// Обработчик для завершения задачи
func DoneTask(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		task := models.Task{}
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		tx := db.WithContext(ctx).Where("id = ?", id).First(&task)
		if tx.Error != nil {

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.SaveTaskError,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)

			return
		}

		if task.Repeat == "" {
			tx := db.Delete(&task)

			if tx.Error != nil {

				errorText, _ := json.Marshal(models.ResponseError{
					MyError: models.DeleteTaskError,
				})

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				w.Write(errorText)

				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
			return
		}

		newDate, err := utils.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			errorText, _ := json.Marshal(err)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)

			return
		}

		task.Date = newDate

		tx = db.Model(&task).Updates(&task)
		if tx.Error != nil {

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.SaveTaskError,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	}
}
