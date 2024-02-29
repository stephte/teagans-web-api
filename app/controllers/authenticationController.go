package controllers

import (
	"teagans-web-api/app/utilities/httpUtils"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/services"
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var dto dtos.LoginDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.LoginService{BaseService: baseService}
	tokenDTO, maxAge, errDTO := service.LoginUser(dto, true)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer: %s", tokenDTO.Token))
	w.Header().Set("X-CSRF-Token", tokenDTO.CSRF)
	w.Header().Set("Expires", strconv.FormatInt(maxAge, 10))

	httpUtils.SetAuthCookie(w, tokenDTO.Token, maxAge, false)

	w.WriteHeader(http.StatusNoContent)
}


func Logout(w http.ResponseWriter, r *http.Request) {
	httpUtils.DeleteAuthCookie(w, false)

	w.WriteHeader(http.StatusNoContent)
}


func StartPWReset(w http.ResponseWriter, r *http.Request) {
	var dto dtos.EmailDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.LoginService{BaseService: baseService}

	errDTO := service.StartPWReset(dto)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, map[string]string{"msg": "Password reset email will be sent if a user with that email exists."}, 200)
}


func ConfirmPasswordResetToken(w http.ResponseWriter, r *http.Request) {
	var dto dtos.ConfirmResetTokenDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.LoginService{BaseService: baseService}

	tokenDTO, maxAge, errDTO := service.ConfirmResetToken(dto)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer: %s", tokenDTO.Token))
	w.Header().Set("X-CSRF-Token", tokenDTO.CSRF)
	w.Header().Set("Expires", strconv.FormatInt(maxAge, 10))

	httpUtils.SetAuthCookie(w, tokenDTO.Token, maxAge, true)

	w.WriteHeader(http.StatusNoContent)
}


func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var dto dtos.ResetPWDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.LoginService{BaseService: baseService}
	tokenDTO, maxAge, errDTO := service.UpdateUserPassword(dto)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer: %s", tokenDTO.Token))
	w.Header().Set("X-CSRF-Token", tokenDTO.CSRF)
	w.Header().Set("Expires", strconv.FormatInt(maxAge, 10))

	httpUtils.DeleteAuthCookie(w, true)
	httpUtils.SetAuthCookie(w, tokenDTO.Token, maxAge, false)

	w.WriteHeader(http.StatusNoContent)
}
