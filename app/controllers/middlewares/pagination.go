package middlewares

import (
	"chi-users-project/app/utilities/http_utils"
	"chi-users-project/app/services/dtos"
	"net/http"
	"strconv"
	"context"
)


func GetPaginationDTO(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sort := r.URL.Query().Get("order")
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		page, limit := 0, 0
		var err error

		if pageStr != "" {
			page, err = strconv.Atoi(pageStr)
			if err != nil {
				errDTO := dtos.CreateErrorDTO(err, 400, false)
				http_utils.RenderErrorJSON(w, r, errDTO)
				return
			}
		}

		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				errDTO := dtos.CreateErrorDTO(err, 400, false)
				http_utils.RenderErrorJSON(w, r, errDTO)
				return
			}
		}

		pageDTO := dtos.PaginationDTO{}
		pageDTO.Init(sort, page, limit)

		ctx := context.WithValue(r.Context(), "paginationDTO", pageDTO)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
