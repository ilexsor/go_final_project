package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ilexsor/internal/models"
	"gorm.io/gorm"
)

// Обработчик для GET /api/task/?id={id}
func GetTask(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		task := models.Task{}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		id := r.URL.Query().Get("id")

		tx := db.WithContext(ctx).Where("id = ?", id).First(&task)

		if tx.Error != nil {
			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.GetTaskError,
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)
			return
		}

		resp, err := json.Marshal(task)

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
