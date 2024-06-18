package app

import (
	"archive/zip"
	"github.com/restartfu/decryptmypack/app/minecraft"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	commonServers = []string{
		"zeqa.net",
		"play.galaxite.net",
	}
	downloading = sync.Map{}
)

func init() {
	for _, server := range commonServers {
		go periodicallyDownloadPacks(server)
	}
}

func periodicallyDownloadPacks(server string) {
	for {
		time.Sleep(time.Minute)
		if err := downloadPacksFromServer(server); err != nil {
			// Log the error (could use a proper logging framework)
			continue
		}
		time.Sleep(time.Duration(60/len(commonServers)) * time.Minute)
	}
}

func downloadPacksFromServer(server string) error {
	conn, err := minecraft.Connect(server)
	if err != nil {
		return err
	}
	defer conn.Close()

	packs := conn.ResourcePacks()
	if len(packs) == 0 {
		return nil
	}

	if err := os.MkdirAll("packs/"+server, 0777); err != nil {
		return err
	}

	filePath := "packs/" + server + "/" + server + ".zip"
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	zipFile := zip.NewWriter(f)
	defer zipFile.Close()

	for _, pack := range packs {
		buf, err := minecraft.EncodePack(pack)
		if err != nil {
			return err
		}
		if pack.Encrypted() {
			buf, err = minecraft.DecryptPack(buf, pack.ContentKey())
			if err != nil {
				return err
			}
		}

		p, err := zipFile.Create(pack.Name() + ".zip")
		if err != nil {
			return err
		}
		if _, err = p.Write(buf); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) download(w http.ResponseWriter, r *http.Request) {
	target := r.FormValue("target")
	if target == "" {
		http.Error(w, "missing target", http.StatusBadRequest)
		return
	}

	target = strings.Split(target, ":")[0]

	if c, ok := downloading.Load(target); ok {
		<-c.(chan struct{})
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")

	filePath := "packs/" + target + "/" + target + ".zip"
	if fileExistsAndFresh(filePath, time.Minute*60) {
		serveFile(w, r, filePath)
		return
	}

	c := make(chan struct{})
	downloading.Store(target, c)
	defer func() {
		close(c)
		downloading.Delete(target)
	}()

	if err := downloadPacksFromServer(target); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	serveFile(w, r, filePath)
}

func fileExistsAndFresh(filePath string, maxAge time.Duration) bool {
	if fi, err := os.Stat(filePath); err == nil {
		return time.Since(fi.ModTime()) <= maxAge
	}
	return false
}

func serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+strings.Split(filePath, "/")[1]+".zip\"")
	http.ServeFile(w, r, filePath)
}
