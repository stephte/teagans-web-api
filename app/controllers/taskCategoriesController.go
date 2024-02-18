package controllers

import (
	"teagans-web-api/app/utilities/httpUtils"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"encoding/json"
	"net/http"
)

func CreateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var dto dtos.TaskCategoryInDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 0, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.TaskCategoryService{BaseService: baseService}

	tcDTO, errDTO := service.CreateTaskCategory(dto)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, tcDTO)
	w.WriteHeader(http.StatusCreated)
}

func UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "categoryId")

	var dto dtos.TaskCategoryInDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 0, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.TaskCategoryService{BaseService: baseService}

	tcDTO, errDTO := service.UpdateTaskCategory(dto, categoryIdStr)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, tcDTO)
}

func DeleteTaskCategory(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "categoryId")

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.TaskCategoryService{BaseService: baseService}

	errDTO := service.DeleteTaskCategory(categoryIdStr)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.NoContent(w, r)
}

func GetTaskCategoryTasks(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "categoryId")

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.TaskCategoryService{BaseService: baseService}

	taskListDTO, errDTO := service.GetTaskCategoryTasks(categoryIdStr)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	render.JSON(w, r, taskListDTO)
}
