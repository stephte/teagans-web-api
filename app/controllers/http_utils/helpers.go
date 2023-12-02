package http_utils

import (
	"youtube-downloader/app/services/dtos"
	"github.com/go-chi/render"
	"net/http"
	"strings"
)

func RenderErrorJSON(w http.ResponseWriter, r *http.Request, errorDTO dtos.ErrorDTO) {
	if errorDTO.Status == 0 {
		errorDTO.Status = 400
	}

	w.WriteHeader(errorDTO.Status)

	render.JSON(w, r, errorDTO)
}


func BlankSuccessResponse(w http.ResponseWriter, r *http.Request) {
	// is there a better way to return just a blank 200 response?
	w.WriteHeader(200)
	render.PlainText(w, r, "")
}


func GetRequestPath(r *http.Request) string {
	url := r.URL.String()

	// now strip out everything after the '?' (if any)
	stringsArr := strings.Split(url, "?")

	return stringsArr[0]
}
