package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ilexsor/internal/models"
	"github.com/ilexsor/internal/utils"
	"gorm.io/gorm"
)

// Обработчик для Delete /task
func DeleteTask(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Query().Get("id")
		if !utils.CheckId(id) {

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.IdFormatError,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)
			return
		}

		tx := db.Delete(&models.Task{}, id)
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
