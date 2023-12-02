package controllers

import (
	"youtube-downloader/app/controllers/http_utils"
	"youtube-downloader/app/services/dtos"
	"youtube-downloader/app/services"
	"github.com/go-chi/render"
	"encoding/json"
	"net/http"
)


func Login(w http.ResponseWriter, r *http.Request) {
	var dto dtos.LoginDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		http_utils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}
	
	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.LoginService{BaseService: baseService}

	tokenDTO, errDTO := service.LoginUser(dto, true)

	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, tokenDTO)
}


func StartPWReset(w http.ResponseWriter, r *http.Request) {
	var dto dtos.EmailDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		http_utils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.LoginService{BaseService: baseService}

	errDTO := service.StartPWReset(dto)

	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, map[string]string{"msg": "Password reset email will be sent if a user with that email exists."})
}


func ConfirmPasswordResetToken(w http.ResponseWriter, r *http.Request) {
	var dto dtos.ConfirmResetTokenDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		http_utils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.LoginService{BaseService: baseService}

	res, errDTO := service.ConfirmResetToken(dto)

	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, res)
}


func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var dto dtos.ResetPWDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		http_utils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.LoginService{BaseService: baseService}
	
	tokenDTO, errDTO := service.UpdateUserPassword(dto)

	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
	}

	render.JSON(w, r, tokenDTO)
}
