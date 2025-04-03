package handlers

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request){
	r.Header.Set("Content-Type", "text/html")
	w.WriteHeader(http.Status.OK)
}