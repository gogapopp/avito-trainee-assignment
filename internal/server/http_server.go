package httpserver

import (
	"assignment/internal/config"
	httphandlers "assignment/internal/server/http_handlers"
	httpmiddleware "assignment/internal/server/http_middlewares"
	"assignment/internal/service"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func New(cfg *config.Config, logger *zap.SugaredLogger, service *service.Service) *http.Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/dummyLogin", httphandlers.DummyLoginHandler(logger, service))
	r.Post("/register", httphandlers.RegisterHandler(logger, service))
	r.Post("/login", httphandlers.LoginHandler(logger, service))

	r.Group(func(r chi.Router) {
		r.Use(httpmiddleware.Auth(cfg.JWTSecret))

		r.Route("/pvz", func(r chi.Router) {
			r.With(httpmiddleware.RequireRole("moderator")).Post("/", httphandlers.CreatePVZHandler(logger, service))
			r.Get("/", httphandlers.GetPVZListHandler(logger, service))
			r.With(httpmiddleware.RequireRole("client")).Post("/{pvzId}/close_last_reception", httphandlers.CloseLastReceptionHandler(logger, service))
			r.With(httpmiddleware.RequireRole("client")).Post("/{pvzId}/delete_last_product", httphandlers.DeleteLastProductHandler(logger, service))
		})

		r.Route("/receptions", func(r chi.Router) {
			r.With(httpmiddleware.RequireRole("client")).Post("/", httphandlers.CreateReceptionHandler(logger, service))
		})

		r.Route("/products", func(r chi.Router) {
			r.With(httpmiddleware.RequireRole("client")).Post("/", httphandlers.AddProductHandler(logger, service))
		})
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", cfg.HTTPServer.Port),
		Handler:           r,
		ReadHeaderTimeout: cfg.HTTPServer.ReadTimeout,
		WriteTimeout:      cfg.HTTPServer.WriteTimeout,
	}

	return server
}
