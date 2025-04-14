package httphandlers

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

type receptionService interface {
	CreateReception(ctx context.Context, req models.NewReceptionRequest) (models.Reception, error)
	CloseLastReception(ctx context.Context, pvzID string) (models.Reception, error)
}

func CreateReceptionHandler(logger *zap.SugaredLogger, service receptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http_handlers.reception.CreateReceptionHandler"
		ctx := r.Context()
		logger := logger.With("req_id", middleware.GetReqID(ctx))

		var req models.NewReceptionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusBadRequest, "invalid request body")
			return
		}

		reception, err := service.CreateReception(ctx, req)
		if err != nil {
			var valErr validator.ValidationErrors
			logger.Errorf("%s: %w", op, err)

			switch {
			case errors.Is(err, repository.ErrPVZNotFound):
				errorJSONResponse(w, http.StatusBadRequest, "pvz not found")
			case errors.Is(err, repository.ErrActiveReceptionExist):
				errorJSONResponse(w, http.StatusBadRequest, "active reception already exists")
			case errors.As(err, &valErr):
				errorJSONResponse(w, http.StatusBadRequest, "validation error")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		jsonResponse(w, http.StatusCreated, reception)
	}
}

func CloseLastReceptionHandler(logger *zap.SugaredLogger, service receptionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http_handlers.reception.CloseLastReceptionHandler"
		ctx := r.Context()
		logger := logger.With("req_id", middleware.GetReqID(ctx))

		pvzID := chi.URLParam(r, "pvzId")
		if pvzID == "" {
			errorJSONResponse(w, http.StatusBadRequest, "pvz id is required")
			return
		}

		reception, err := service.CloseLastReception(ctx, pvzID)
		if err != nil {
			logger.Errorf("%s: %w", op, err)

			switch {
			case errors.Is(err, repository.ErrNoActiveReception):
				errorJSONResponse(w, http.StatusBadRequest, "no active reception found")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		jsonResponse(w, http.StatusOK, reception)
	}
}
