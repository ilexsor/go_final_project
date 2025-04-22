package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ilexsor/internal/models"
	log "github.com/sirupsen/logrus"
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

		search := r.URL.Query().Get("search")

		// Если передан параметр search, то необходимо проверть его на условия
		// Если это дата, то преобразовать ее в нужный формат и выполнять поиск по полую date
		// Если это строка, то отправить запрос в БД на поиск записей с такой строкой
		// Если условия search не задан, то вернуть все таски
		if search != "" || len(search) != 0 {

			//Пробуем распарсить переданную дату
			date, err := time.Parse("02.01.2006", search)

			// Если дата не парсится, то делаюм заключение, что нужен текстовый поиск
			// Логируем пробему при парсинге и пробуем найти по полю title
			if err != nil {
				errorText, _ := json.Marshal(models.ResponseError{
					MyError: models.DateError,
				})

				log.WithFields(log.Fields{
					"serach": date,
					"msg":    "parse date error, try to find plain text",
				}).Errorf(string(errorText))

				// Поиск задач по полю title
				tx := db.WithContext(ctx).Where("title LIKE ?", "%"+search+"%").Find(&task)

				if tx.Error != nil {

					errorText, _ := json.Marshal(models.ResponseError{
						MyError: models.GetTasksError,
					})
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusBadRequest)
					w.Write(errorText)
					return
				}

				// Если задачи по полю title не найдены, то пробуем искать по comment
				if len(task) == 0 {
					tx := db.WithContext(ctx).Where("comment LIKE ?", "%"+search+"%").Find(&task)

					if tx.Error != nil {

						errorText, _ := json.Marshal(models.ResponseError{
							MyError: models.GetTasksError,
						})
						w.Header().Set("Content-Type", "application/json")
						w.WriteHeader(http.StatusBadRequest)
						w.Write(errorText)
						return
					}
					// Если задачи по comment не найдены, то отдаем пустой ответ
					if len(task) == 0 {
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
						return
					}
					// Если задачи по comment найдены, отдаем их
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
					return
				}

				// Если найдены задачи по title, то формируем ответ
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
				return
			}

			// Если дату удалось распарсить, ищем по ней
			formattedDate := date.Format("20060102")

			tx := db.WithContext(ctx).Where("date = ?", formattedDate).Find(&task)

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
			return
		}

		// Если не передан параметр search забираем все задачи из БД

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
