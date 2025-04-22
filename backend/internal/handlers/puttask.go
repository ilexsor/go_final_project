package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ilexsor/internal/models"
	"github.com/ilexsor/internal/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Обработчик для Put /task
func PutTask(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		task := &models.Task{}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			log.WithFields(log.Fields{
				"readBody": err,
			}).Error()

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.RequestBodyError,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)

			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.UnmarshalError,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)

			return
		}

		task, err = utils.CheckTask(task)

		if err != nil {

			errorText, _ := json.Marshal(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)

			return
		}

		if !utils.CheckId(task.ID) {

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.IdFormatError,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)
			return

		}

		tx := db.WithContext(ctx).Save(&task)
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
