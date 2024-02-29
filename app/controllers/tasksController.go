package controllers

import (
	"teagans-web-api/app/utilities/httpUtils"
	"teagans-web-api/app/services/dtos"
	"teagans-web-api/app/services"
	"github.com/go-chi/chi/v5"
	"encoding/json"
	"net/http"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var dto dtos.TaskInDTO
	bindErr := json.NewDecoder(r.Body).Decode(&dto)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 0, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.TaskService{BaseService: baseService}

	taskDTO, errDTO := service.CreateTask(dto)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, taskDTO, 201)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	taskIdStr := chi.URLParam(r, "taskId")

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.TaskService{BaseService: baseService}

	taskOutDTO, errDTO := service.GetTask(taskIdStr)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, taskOutDTO, 200)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskIdStr := chi.URLParam(r, "taskId")

	var data map[string]interface{}
	bindErr := json.NewDecoder(r.Body).Decode(&data)
	if bindErr != nil {
		httpUtils.RenderErrorJSON(w, r, dtos.CreateErrorDTO(bindErr, 0, false))
		return
	}

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.TaskService{BaseService: baseService}

	taskOutDTO, errDTO := service.UpdateTask(data, taskIdStr)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	httpUtils.RenderJSON(w, taskOutDTO, 200)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskIdStr := chi.URLParam(r, "taskId")

	baseService := r.Context().Value("BaseService").(*services.BaseService)
	service := services.TaskService{BaseService: baseService}

	errDTO := service.DeleteTask(taskIdStr)
	if errDTO.Exists() {
		httpUtils.RenderErrorJSON(w, r, errDTO)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
