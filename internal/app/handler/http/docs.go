package http

import (
	"net/http"
	"os"

	"web-studio-backend/internal/app/handler/http/httphelp"
)

const apiFilePath = "api/docs/api.html"

func getApiDocs(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat(apiFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		httphelp.SendError(err, w)
		return
	}

	http.ServeFile(w, r, apiFilePath)
}
