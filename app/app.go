package app

import (
	"github.com/gorilla/mux"
	"net/http"
)

var router = mux.NewRouter()

type App struct {
}

func (a *App) ListenAndServe(addr string) error {
	router.HandleFunc("/download", a.download).Queries("target", "{target}")

	router.HandleFunc("/", serveFileFunc("./frontend/static/home.html"))
	router.HandleFunc("/style.css", serveFileFunc("./frontend/static/style.css"))
	router.HandleFunc("/src/script.js", serveFileFunc("./frontend/src/script.js"))
	router.HandleFunc("/assets/quicksand_bold.ttf", serveFileFunc("./frontend/assets/quicksand_bold.ttf"))
	return http.ListenAndServe(addr, router)
}

func serveFileFunc(name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	}
}
