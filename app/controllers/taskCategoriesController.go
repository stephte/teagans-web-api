package controllers

import (
	"teagans-web-api/app/utilities/httpUtils"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/services"
	"github.com/go-chi/chi/v5"
	"encoding/json"
	"net/http"
	"strconv"
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

	httpUtils.RenderJSON(w, tcDTO, 201)
}

func UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "categoryId")

	var data map[string]interface{}
	bindErr := json.NewDecoder(r.Body).Decode(&data)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 0, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.TaskCategoryService{BaseService: baseService}

	tcDTO, errDTO := service.UpdateTaskCategory(data, categoryIdStr)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, tcDTO, 200)
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

	w.WriteHeader(http.StatusNoContent)
}

// add query params to filter by status and if 'cleared' or not
func GetTaskCategoryTasks(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "categoryId")
	statusQuery := r.URL.Query().Get("status")
	getCleared, _ := strconv.ParseBool(r.URL.Query().Get("cleared")) // always assume false if not true, not required that this be passed in

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.TaskCategoryService{BaseService: baseService}

	taskListDTO, errDTO := service.GetTaskCategoryTasks(categoryIdStr, statusQuery, getCleared)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, taskListDTO, 200)
}
