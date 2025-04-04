package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/ilexsor/internal/handlers"
	"github.com/ilexsor/internal/utils"
	log "github.com/sirupsen/logrus"
)

const (
	frontendDir = http.Dir("../../../web")
)

func main() {
	port := utils.GetServerPort()
	log.SetFormatter(&log.JSONFormatter{})

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	handlers.FileServer(router, "/", frontendDir)

	log.WithFields(log.Fields{
		"server status": "starting",
	}).Info("starting on port ", port)

	if err := http.ListenAndServe(port, router); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
