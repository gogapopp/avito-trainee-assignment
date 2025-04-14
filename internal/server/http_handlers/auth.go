package httphandlers

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

type authService interface {
	RegisterUser(ctx context.Context, req models.AuthRequest) (models.User, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
	DummyLogin(ctx context.Context, role string) (string, error)
}

func RegisterHandler(logger *zap.SugaredLogger, service authService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http_handlers.auth.RegisterHandler"
		ctx := r.Context()
		logger := logger.With("req_id", middleware.GetReqID(ctx))

		var req models.AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusBadRequest, "invalid request body")
			return
		}

		user, err := service.RegisterUser(ctx, req)
		if err != nil {
			var valErr validator.ValidationErrors
			logger.Errorf("%s: %w", op, err)

			switch {
			case errors.Is(err, repository.ErrUserExists):
				errorJSONResponse(w, http.StatusBadRequest, "user already exists")
			case errors.As(err, &valErr):
				errorJSONResponse(w, http.StatusBadRequest, "validation error")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		jsonResponse(w, http.StatusCreated, user)
	}
}

func LoginHandler(logger *zap.SugaredLogger, service authService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http_handlers.auth.LoginHandler"
		ctx := r.Context()
		logger := logger.With("req_id", middleware.GetReqID(ctx))

		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusBadRequest, "invalid request body")
			return
		}

		token, err := service.LoginUser(ctx, req.Email, req.Password)
		if err != nil {
			logger.Errorf("%s: %w", op, err)
			switch {
			case errors.Is(err, repository.ErrInvalidCredentials):
				errorJSONResponse(w, http.StatusUnauthorized, "invalid credentials")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		jsonResponse(w, http.StatusOK, token)
	}
}

func DummyLoginHandler(logger *zap.SugaredLogger, service authService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http_handlers.auth.DummyLoginHandler"
		ctx := r.Context()
		logger := logger.With("req_id", middleware.GetReqID(ctx))

		var req struct {
			Role string `json:"role"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Role != "client" && req.Role != "moderator" {
			errorJSONResponse(w, http.StatusBadRequest, "invalid role")
			return
		}

		token, err := service.DummyLogin(ctx, req.Role)
		if err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusInternalServerError, "internal server error")
			return
		}

		jsonResponse(w, http.StatusOK, token)
	}
}
