package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	"github.com/ilexsor/internal/models"
	"github.com/ilexsor/internal/utils"
	log "github.com/sirupsen/logrus"
)

func Signin(w http.ResponseWriter, r *http.Request) {

	pass := os.Getenv("TODO_PASSWORD")

	if pass == "" {
		errorText, _ := json.Marshal(models.ResponseError{
			MyError: models.IncorrectPassword,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errorText)

		return
	}

	token := utils.GetToken(pass)

	var buf bytes.Buffer
	passModel := models.Auth{}

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

	if err = json.Unmarshal(buf.Bytes(), &passModel); err != nil {

		errorText, _ := json.Marshal(models.ResponseError{
			MyError: models.UnmarshalError,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorText)

		return
	}

	if pass != passModel.Password {
		errorText, _ := json.Marshal(models.ResponseError{
			MyError: models.IncorrectPassword,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errorText)

		return
	}

	response, _ := json.Marshal(models.Token{
		Token: token,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
