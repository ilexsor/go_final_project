package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/ilexsor/internal/database"
	"github.com/ilexsor/internal/handlers"
	"github.com/ilexsor/internal/utils"
	log "github.com/sirupsen/logrus"
)

const (
	frontendDir = http.Dir("../../../web")
	dbFile      = "../../internal/database/scheduler.db"
)

func main() {
	// Читаем переменную среды для порта
	port := utils.GetServerPort()

	log.SetFormatter(&log.JSONFormatter{})

	// Инициализируем БД
	_, err := database.NewSqliteDB(dbFile)
	if err != nil {
		log.WithFields(log.Fields{
			"migration": "migration error",
		}).Errorf("error: %v", err)
		return
	}

	// Подключаем CHI router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	
	router.Route("/api", func(router chi.Router){
		router.Get("/nextdate", handlers.NextDayHandler)
	})
	
	
	handlers.FileServer(router, "/", frontendDir)

	log.WithFields(log.Fields{
		"server status": "starting",
	}).Info("starting on port ", port)

	// Запускаем сервер
	if err := http.ListenAndServe(port, router); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
