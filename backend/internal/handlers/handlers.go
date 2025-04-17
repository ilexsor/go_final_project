package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	//"fmt"
	"bytes"
	"net/http"
	"strconv"
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
		nowStr = time.Now().Format(utils.DateFormat)
	}

	if repeat == "" {
		errorText, _ := json.Marshal(models.ResponseError{
			MyError: models.RepeatParamError,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorText)

		log.WithFields(log.Fields{
			"now":    nowStr,
			"date":   dstart,
			"repeat": repeat,
		}).Errorf("request: %v", r.URL.String())
		return
	}

	now, err := time.Parse(utils.DateFormat, nowStr)
	if err != nil {
		errorText, _ := json.Marshal(models.ResponseError{
			MyError: models.DateError,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorText)

		log.WithFields(log.Fields{
			"now":    nowStr,
			"date":   dstart,
			"repeat": repeat,
		}).Errorf("request:%v", r.URL.String())
		return
	}

	if _, err := time.Parse(utils.DateFormat, dstart); err != nil {

		errorText, _ := json.Marshal(models.ResponseError{
			MyError: models.DateError,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorText)

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

// Обработчик для POST /api/task
func AddTask(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		task := &models.Task{}
		taskResponse := &models.TaskResponse{}
		schedule := &models.Scheduler{}

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

		schedule, err = utils.ConvertTaskToSchedule(task, schedule)
		if err != nil {
			log.WithFields(log.Fields{
				"convertTask": "error",
			}).Errorf("error:%v", err)

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.ConvertationError,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)
			return

		}

		tx := db.WithContext(ctx).Create(&schedule)
		if tx.Error != nil {

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.SaveTaskError,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)

			return
		}
		id := strconv.Itoa(schedule.ID)
		taskResponse.ID = id

		response, err := json.Marshal(taskResponse)

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
		w.Write(response)
	}
}

// Обработчик для GET /api/tasks
func GetTasks(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var schedules []models.Schedule
		response := models.TasksResponse{
			Tasks: []models.Schedule{},
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		tx := db.WithContext(ctx).Limit(models.TaskLimit).Find(&schedules)

		if tx.Error != nil {

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.GetTasksError,
			})
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)
			return
		}

		response.Tasks = schedules
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

func PutTask(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		schedule := &models.Scheduler{}
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
		fmt.Println(task)
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

		schedule, err = utils.ConvertTaskToSchedule(task, schedule)
		if err != nil {
			log.WithFields(log.Fields{
				"convertTask": "error",
			}).Errorf("error:%v", err)

			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.ConvertationError,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorText)
			return

		}

		tx := db.WithContext(ctx).Save(&schedule)
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
