package utils

import (
	"github.com/go-chi/chi/v5"
)


func RegisterHandlers(r *chi.Mux) {
	r.Get("/", handlers.IndexHandler)
}