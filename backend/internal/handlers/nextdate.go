package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ilexsor/internal/models"
	"github.com/ilexsor/internal/utils"
	log "github.com/sirupsen/logrus"
)

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
