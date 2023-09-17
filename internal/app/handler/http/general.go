package http

import (
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"web-studio-backend/internal/app/handler/http/httphelp"
)

const (
	apiSwaggerFileName = "swagger.json"
	apiPageFileName    = "apidocs.html"
)

var (
	apiDocsDir = filepath.Join("web", "static", "apidocs")
)

func getApiDocsSwagger(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(apiDocsDir, apiSwaggerFileName)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		httphelp.SendError(err, w)
		return
	}

	http.ServeFile(w, r, filePath)
}

func getApiDocs(w http.ResponseWriter, r *http.Request) {
	filePath := filepath.Join(apiDocsDir, apiPageFileName)

	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		httphelp.SendError(err, w)
		return
	}

	http.ServeFile(w, r, filePath)
}

func getStatic(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/static/") {
		w.WriteHeader(http.StatusNotFound)
		slog.Error("url does not have /static/ prefix, but static handler is accessed")
		return
	}
	filename := strings.TrimPrefix(r.URL.Path, "/static/")

	http.ServeFile(w, r, path.Join("web", "static", filename))
}
