package controllers

import (
	"youtube-downloader/app/utilities/http_utils"
	"youtube-downloader/app/services/dtos"
	"youtube-downloader/app/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"encoding/json"
	"net/http"
)


func UsersIndex(w http.ResponseWriter, r *http.Request) {
	paginationDTO := r.Context().Value("paginationDTO").(dtos.PaginationDTO)

	path := http_utils.GetRequestPath(r)

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	result, errDTO := service.GetUsers(paginationDTO, path)
	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, result)
}


func FindUser(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	userDTO, errDTO := service.GetUser(userIdStr)
	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, userDTO)
}


func CreateUser(w http.ResponseWriter, r *http.Request) {
	var dto dtos.CreateUserDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		http_utils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	userDTO, errDTO := service.CreateUser(dto)

	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, userDTO)
}


// PATCH version of User update
// this endpoint validates the request data against the UserDTO,
// but keeps it as a map so that only the included data is updated
// (GORM only updates non-zero fields when updating with struct)
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")

	var data map[string]interface{}
	bindErr := json.NewDecoder(r.Body).Decode(&data)
	if bindErr != nil {
		http_utils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	userDTO, errDTO := service.UpdateUser(userIdStr, data)

	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, userDTO)
}


// PUT version of User update (expects all user data) (prefer above PATCH version)
func UpdateUserOG(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")

	var dto dtos.UserDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		http_utils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 400, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	userDTO, errDTO := service.UpdateUserOG(userIdStr, dto)

	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, userDTO)
}


func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	errDTO := service.DeleteUser(userIdStr)

	if errDTO.Exists() {
		http_utils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.NoContent(w, r)
}
