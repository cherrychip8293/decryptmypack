package template

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func NewFS(path string, replacer *strings.Replacer) func(http.ResponseWriter, *http.Request) {
	buf, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	content := replacer.Replace(string(buf))

	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = fmt.Fprint(w, content)
	}
}
