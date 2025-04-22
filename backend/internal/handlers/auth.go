package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ilexsor/internal/models"
	"github.com/ilexsor/internal/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		pass := os.Getenv("TODO_PASSWORD")

		if pass == "" || len(pass) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		token := utils.GetToken(pass)

		cookie, err := r.Cookie("token")

		if err != nil {
			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.ReadCookieError,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(errorText)

			return
		}

		if cookie.Value != token {
			errorText, _ := json.Marshal(models.ResponseError{
				MyError: models.AuthRequired,
			})

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(errorText)
			return
		}

		next.ServeHTTP(w, r)

	})
}
