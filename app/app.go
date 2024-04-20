package app

import (
	"github.com/gorilla/mux"
	"github.com/restartfu/decryptmypack/app/template"
	"net/http"
	"strings"
)

var (
	router = mux.NewRouter()
)

type App struct {
}

func (a *App) ListenAndServe(addr string, dev bool) error {
	downloadAddr := "https://decryptmypack.com:443"
	if dev {
		downloadAddr = "http://localhost:8080"
	}

	router.HandleFunc("/download", a.download).Queries("target", "{target}")

	router.HandleFunc("/", serveFileFunc("./frontend/static/home.html"))
	router.HandleFunc("/style.css", serveFileFunc("./frontend/static/style.css"))
	router.HandleFunc("/src/script.js", template.NewFS("./frontend/src/script.js", strings.NewReplacer(
		"$DOWNLOAD_ADDR", downloadAddr,
	)))
	router.HandleFunc("/assets/Quicksand_Bold.otf", serveFileFunc("./frontend/assets/Quicksand_Bold.otf"))

	if dev {
		return http.ListenAndServe(addr, router)
	}
	return http.ListenAndServeTLS(addr, "./certificate.crt", "./private.key", router)
}

func serveFileFunc(name string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, name)
	}
}
