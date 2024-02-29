package controllers

import (
	"teagans-web-api/app/utilities/httpUtils"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/services"
	"github.com/go-chi/chi/v5"
	"encoding/json"
	"net/http"
)


func UsersIndex(w http.ResponseWriter, r *http.Request) {
	paginationDTO := r.Context().Value("paginationDTO").(dtos.PaginationDTO)

	path := httpUtils.GetRequestPath(r)

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	result, errDTO := service.GetUsers(paginationDTO, path)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, result, 200)
}


func GetUser(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	userDTO, errDTO := service.GetUser(userIdStr)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, userDTO, 200)
}


func CreateUser(w http.ResponseWriter, r *http.Request) {
	var dto dtos.CreateUserDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 0, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	userDTO, errDTO := service.CreateUser(dto)

	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, userDTO, 201)
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
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 0, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	userDTO, errDTO := service.UpdateUser(userIdStr, data)

	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, userDTO, 200)
}


// PUT version of User update (expects all user data) (prefer above PATCH version)
func UpdateUserOG(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")

	var dto dtos.UserInDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 0, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	userDTO, errDTO := service.UpdateUserOG(userIdStr, dto)

	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, userDTO, 200)
}


func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	errDTO := service.DeleteUser(userIdStr)

	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func GetUserTaskCategories(w http.ResponseWriter, r *http.Request) {
	userIdStr := chi.URLParam(r, "userId")

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.UserService{BaseService: baseService}

	categoriesDTO, errDTO := service.GetUserTaskCategories(userIdStr)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, categoriesDTO, 200)
}
