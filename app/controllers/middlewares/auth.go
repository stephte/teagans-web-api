package middlewares

import (
	"youtube-downloader/app/utilities/http_utils"
	"youtube-downloader/app/services/dtos"
	"youtube-downloader/app/services"
	"net/http"
	"context"
	"strings"
	"errors"
)


func ValidateJWT(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, authService := getAuthTokenAndService(r, false)
		csrf := r.Header.Get("X-CSRF-Token")

		jwtValid, tokenErrDTO := authService.ValidateJWT(token, csrf, false)

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
		token, authService := getAuthTokenAndService(r, true)
		csrf := r.Header.Get("X-CSRF-Token")

		jwtValid, tokenErrDTO := authService.ValidateJWT(token, csrf, true)

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
		token, authService := getAuthTokenAndService(r, false)
		csrf := r.Header.Get("X-CSRF-Token")

		jwtValid, tokenErrDTO := authService.ValidateJWT(token, csrf, false)

		ctx, _ := handleErrDTO(jwtValid, tokenErrDTO, r)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AllowLocalHost(next http.Handler) (http.Handler) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Content-Disposition, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Expose-Headers", "Content-Type, Content-Disposition, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
    	w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    	w.Header().Add("Access-Control-Allow-Credentials", "true")

    	if r.Method == "OPTIONS" {
	        http.Error(w, "No Content", http.StatusNoContent)
	        return
	    }

		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}

// ----- Private -----

func getAuthTokenAndService(r *http.Request, pwReset bool) (string, services.AuthService) {
	var token string
	jwt := r.Header.Get("Authorization")

	// if auth token not in headers, get it from cookies
	if jwt == "" {
		authCookie, noCookieErr := http_utils.GetAuthCookie(r, pwReset)
		if noCookieErr == nil {
			token = authCookie.Value
		}
	} else {
		token = strings.Replace(jwt, "Bearer ", "", 1)
	}

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
