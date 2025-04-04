package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
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
