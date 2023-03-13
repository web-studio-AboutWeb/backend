package http

import (
	"net/http"
	"os"
)

const apiFilePath = "api/docs/api.html"

func (s *server) GetApiDocs(w http.ResponseWriter, r *http.Request) {
	_, err := os.Stat(apiFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		s.sendError(err, w)
		return
	}

	http.ServeFile(w, r, apiFilePath)
}
