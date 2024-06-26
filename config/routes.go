package config

import (
	middle "teagans-web-api/app/controllers/middlewares"
	"github.com/go-chi/chi/v5/middleware"
	"teagans-web-api/app/controllers"
	"teagans-web-api/app/services"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"os/signal"
	"net/http"
	"context"
	"syscall"
	"errors"
	"time"
	"fmt"
	"os"
)

type Router struct {
	Router			*chi.Mux
	logger			zerolog.Logger
	env				string
}

// -------- Routes go here --------

func(this *Router) defineRoutes() {
	r := this.Router

	this.logger.Debug().Msg("Setting up routes")

	// uses CHI_YT_WEBAPP_ORIGIN to allow requests from given origin
	r.Use(middle.SetCORS)

	r.Get("/", controllers.Index)

	// authentication + password handling
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", controllers.Login)
		r.Post("/logout", controllers.Logout)
		r.Post("/reset-password", controllers.StartPWReset)
		r.Post("/confirm-reset-token", controllers.ConfirmPasswordResetToken)
		r.With(middle.ValidatePWResetJWT).Post("/update-password", controllers.UpdatePassword)
	})

	// users
	r.Route("/users", func(r chi.Router) {
		r.With(middle.ValidateJWT).With(middle.GetPaginationDTO).Get("/", controllers.UsersIndex)
		r.With(middle.ValidateOptionalJWT).Post("/", controllers.CreateUser)

		r.Route("/{userId}", func(r chi.Router) {
			r.Use(middle.ValidateJWT)

			r.Get("/", controllers.GetUser)
			r.Patch("/", controllers.UpdateUser)
			// r.Put("/", controllers.UpdateUserOG)
			r.Delete("/", controllers.DeleteUser)

			r.Get("/task-categories", controllers.GetUserTaskCategories)
		})
	})

	r.Route("/download", func(r chi.Router) {
		r.Use(middle.ValidateJWT)	
		r.Get("/", controllers.DownloadVideo)
	})

	r.Route("/task-categories", func(r chi.Router) {
		r.Use(middle.ValidateJWT)

		r.Post("/", controllers.CreateTaskCategory)
		r.Post("/update-categories", controllers.UpdateTaskCategories)

		r.Route("/{categoryId}", func(r chi.Router) {
			r.Patch("/", controllers.UpdateTaskCategory)
			r.Delete("/", controllers.DeleteTaskCategory)

			r.Get("/tasks", controllers.GetTaskCategoryTasks)
		})
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Use(middle.ValidateJWT)

		r.Post("/", controllers.CreateTask)
		r.Post("/update-tasks", controllers.UpdateTasks)

		r.Route("/{taskId}", func(r chi.Router) {
			r.Get("/", controllers.GetTask)
			r.Patch("/", controllers.UpdateTask)
			r.Delete("/", controllers.DeleteTask)
		})
	})
}

// ----- Other Router/server methods -----

func(this *Router) StartGracefulServer(baseUrl, port string) {
	if this.Router == nil {
		err := errors.New("Router not set up; call SetupRouter before starting the server.")
		this.logger.Error().Err(err).Msg("")
		return
	}

	addr := fmt.Sprintf("%s:%s", baseUrl, port)
	// The HTTP Server
	server := &http.Server{Addr: addr, Handler: this.Router}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		fmt.Print("\nShutting down gracefully...\n\n")

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel() // idk if this is necessary, but got tired of the warning

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				this.logger.Fatal().Msg("Graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			this.logger.Fatal().Err(err).Msg("")
		}
		serverStopCtx()
	}()

	// Run the server
	fmt.Printf("\nServer listening at: %s\n\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		this.logger.Fatal().Err(err).Msg("")
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}


func SetupRouter(logger zerolog.Logger, db *gorm.DB, env string) Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// defined here since it needs access to the database connection
	serviceMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// timeoutContext, _ := context.WithTimeout(context.Background(), 5*time.Second)
			service := services.InitService(logger, db)//.WithContext(timeoutContext))
			ctx := context.WithValue(r.Context(), "BaseService", &service)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
	
	r.Use(serviceMiddleware)

	Router := Router{
		Router: r,
		logger: logger,
		env: env,
	}

	Router.defineRoutes()

	return Router
}

// ----- No longer used -----

func(this *Router) StartServer(port string) error {
	if this.Router == nil {
		return errors.New("Router not set up; call SetupRouter before starting the server.")
	}

	listen_port := fmt.Sprintf(":%s", port)
	fmt.Printf("\nServer listening on localhost:%s\n", port)
	return http.ListenAndServe(listen_port, this.Router)
}


func(this Router) GetHandler() http.Handler {
	return this.Router
}
