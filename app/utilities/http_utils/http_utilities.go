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


func GetRequestPath(r *http.Request) string {
	url := r.URL.String()

	// now strip out everything after the '?' (if any)
	stringsArr := strings.Split(url, "?")

	return stringsArr[0]
}

// sets the Auth cookie on response writer
func SetAuthCookie(w http.ResponseWriter, token string, maxAge int64, pwReset bool) {
	cookie := &http.Cookie{
		Name: getAuthCookieName(pwReset),
		Value: token,
		MaxAge: int(maxAge),
		Domain: "teaganswebapp.com",
		Secure: true, // figure out how to use https for localhost
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path: "/",
	}
	http.SetCookie(w, cookie)

	cookie2 := &http.Cookie{
		Name: "auth-cookie",
		Value: "abc123",
		MaxAge: int(maxAge),
		Domain: "teaganswebapp.com",
		// Secure: true, // figure out how to use https for localhost
		// HttpOnly: true,
		// SameSite: http.SameSiteStrictMode,
		Path: "/",
	}

	http.SetCookie(w, cookie2)

}

// deletes the Auth cookie from response writer
func DeleteAuthCookie(w http.ResponseWriter, pwReset bool) {
	cookie := &http.Cookie{
		Name: getAuthCookieName(pwReset),
		Value: "",
		MaxAge: -1,
		HttpOnly: true,
		Secure: true, // figure out how to use https for localhost
		SameSite: http.SameSiteStrictMode,
		Path: "/",
	}
	http.SetCookie(w, cookie)
}

func GetAuthCookie(r *http.Request, pwReset bool) (*http.Cookie, error) {
	return r.Cookie(getAuthCookieName(pwReset))
}

func getAuthCookieName(pwReset bool) string {
	if pwReset {
		return "yt-downloader-reset"
	}
	return "yt-downloader-auth"
}
