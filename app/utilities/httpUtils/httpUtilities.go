package httpUtils

import (
	"teagans-web-api/app/services/dtos"
	"encoding/json"
	"net/http"
	"strings"
)

func RenderJSON(w http.ResponseWriter, data interface{}, status int) {
	rvJson, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
	}

	// default status to 200
	if status == 0 {
		status = 200
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(rvJson)
}

func RenderErrorJSON(w http.ResponseWriter, r *http.Request, errorDTO dtos.ErrorDTO) {
	if errorDTO.Status == 0 {
		errorDTO.Status = 400
	}

	// w.WriteHeader(errorDTO.Status)
	RenderJSON(w, errorDTO, errorDTO.Status)
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
		Secure: true, // figure out how to use https for localhost
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path: "/",
	}
	http.SetCookie(w, cookie)
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
		return "teagans-app-reset"
	}
	return "teagans-app-auth"
}
