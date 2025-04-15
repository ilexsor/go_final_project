package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/ilexsor/internal/models"
	"github.com/ilexsor/internal/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// FileServer Handler для обработки статических файлов
// Принимает в качестве параметров роутер, путь и каталог со статическими файлами
func FileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, fs.ServeHTTP)
}

// nextDayHandler Обработчик для api/nextdate
func NextDayHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	nowStr := query.Get("now")
	dstart := query.Get("date")
	repeat := query.Get("repeat")

	if nowStr == "" {
		nowStr = time.Now().Format("20060102")
	}

	if repeat == "" {
		http.Error(w, "Parameter 'repeat' are required", http.StatusBadRequest)
		log.WithFields(log.Fields{
			"now":    nowStr,
			"date":   dstart,
			"repeat": repeat,
		}).Errorf("request: %v", r.URL.String())
		return
	}

	now, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Invalid 'now' date format (expected YYYYMMDD)", http.StatusBadRequest)
		log.WithFields(log.Fields{
			"now":    nowStr,
			"date":   dstart,
			"repeat": repeat,
		}).Errorf("request:%v", r.URL.String())
		return
	}

	if _, err := time.Parse("20060102", dstart); err != nil {
		http.Error(w, "Invalid 'dstart' date format", http.StatusBadRequest)
		log.WithFields(log.Fields{
			"now":    nowStr,
			"date":   dstart,
			"repeat": repeat,
		}).Errorf("request:%v", r.URL.String())
		return
	}

	nextDate, err := utils.NextDate(now, dstart, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(nextDate))
}

// Обработчик для /api/task
func AddTask(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		log.WithFields(log.Fields{
			"addTask": "incoming request",
		}).Infof("request:%v", r.URL.String())

		task := models.Scheduler{}
		body, err := io.ReadAll(r.Body)
		log.WithFields(log.Fields{
			"body": body,
		}).Error()

		if err != nil {
			log.WithFields(log.Fields{
				"readBody": "error",
			}).Errorf("request:%v", r.URL.String())
		}

		err = json.Unmarshal(body, &task)
		if err != nil {
			log.WithFields(log.Fields{
				"unmarshal": "error",
			}).Errorf("error:%v", err)
			http.Error(w, "Ошибка анмаршаллинга", http.StatusBadRequest)
			return
		}

		tx := db.WithContext(ctx).Create(&task)
		if tx.Error != nil {
			http.Error(w, "Ошибка при сохранении задачи", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(fmt.Sprintf("id:\"%v\"", task.ID))
		if err != nil {
			http.Error(w, "Ошибка анмаршаллинга", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
