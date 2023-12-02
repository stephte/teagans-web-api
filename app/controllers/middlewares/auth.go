package middlewares

import (
	"youtube-downloader/app/controllers/http_utils"
	"youtube-downloader/app/services/dtos"
	"youtube-downloader/app/services"
	"net/http"
	"context"
	"strings"
	"errors"
)


func ValidateJWT(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, authService := getAuthTokenAndService(r)
		jwtValid, tokenErrDTO := authService.ValidateJWT(token, false)

		ctx, errDTO := handleErrDTO(jwtValid, tokenErrDTO, r)
		if errDTO.Exists() {
			http_utils.RenderErrorJSON(w, r, errDTO)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func ValidatePWResetJWT(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, authService := getAuthTokenAndService(r)

		jwtValid, tokenErrDTO := authService.ValidateJWT(token, true)

		ctx, errDTO := handleErrDTO(jwtValid, tokenErrDTO, r)
		if errDTO.Exists() {
			http_utils.RenderErrorJSON(w, r, errDTO)
			return
		}
	
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


func ValidateOptionalJWT(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, authService := getAuthTokenAndService(r)

		jwtValid, tokenErrDTO := authService.ValidateJWT(token, false)

		ctx, _ := handleErrDTO(jwtValid, tokenErrDTO, r)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// ----- Private -----

func getAuthTokenAndService(r *http.Request) (string, services.AuthService) {
	jwt := r.Header.Get("Authorization")
	token := strings.Replace(jwt, "Bearer ", "", 1)

	service := r.Context().Value("BaseService").(*services.BaseService)
	authService := services.AuthService{BaseService: service}

	return token, authService
}


func handleErrDTO(jwtValid bool, tokenErrDTO dtos.ErrorDTO, r *http.Request) (context.Context, dtos.ErrorDTO) {
	ctx := r.Context()
	var errDTO dtos.ErrorDTO

	if tokenErrDTO.Exists() {
		errDTO = tokenErrDTO
	} else if !jwtValid {
		errDTO = dtos.CreateErrorDTO(errors.New("Token Expired"), 401, true)
	} else {
		ctx = context.WithValue(ctx, "jwtError", false)
	}

	return ctx, errDTO
}
