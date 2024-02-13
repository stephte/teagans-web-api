package controllers

import (
	// "teagans-web-api/app/utilities/http_utils"
	// "teagans-web-api/app/services/dtos"
	// "teagans-web-api/app/services"
	// "github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	// "encoding/json"
	"net/http"
)

func TaskCategoryTasks(w http.ResponseWriter, r *http.Request) {
	// catIdStr := chi.URLParam(r, "categoryId")

	// baseService := r.Context().Value("BaseService").(*services.BaseService)
	// service := services.UserService{BaseService: baseService}

	render.NoContent(w, r)
}
