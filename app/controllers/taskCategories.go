package controllers

import (
	// "teagans-web-api/app/utilities/httpUtils"
	// "teagans-web-api/app/services/dtos"
	// "teagans-web-api/app/services"
	// "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	// "encoding/json"
	"net/http"
)

func CreateTaskCategory(w http.ResponseWriter, r *http.Request) {
	render.NoContent(w, r)
}

func UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	render.NoContent(w, r)
}

func DeleteTaskCategory(w http.ResponseWriter, r *http.Request) {
	render.NoContent(w, r)
}

func GetTaskCategoryTasks(w http.ResponseWriter, r *http.Request) {
	// catIdStr := chi.URLParam(r, "categoryId")

	// baseService := r.Context().Value("BaseService").(*services.BaseService)
	// service := services.TaskCategoryService{BaseService: baseService}

	render.NoContent(w, r)
}
