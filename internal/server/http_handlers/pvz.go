package httphandlers

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

type pvzService interface {
	CreatePVZ(ctx context.Context, pvz models.PVZ) (models.PVZ, error)
	GetPVZList(ctx context.Context, startDate, endDate *string, page, limit int) ([]models.PVZWithReceptions, error)
}

func CreatePVZHandler(logger *zap.SugaredLogger, service pvzService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http_handlers.pvz.CreatePVZHandler"
		ctx := r.Context()
		logger := logger.With("req_id", middleware.GetReqID(ctx))

		var req models.PVZ
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusBadRequest, "invalid request body")
			return
		}

		pvz, err := service.CreatePVZ(ctx, req)
		if err != nil {
			var valErr validator.ValidationErrors
			logger.Errorf("%s: %w", op, err)

			switch {
			case errors.Is(err, repository.ErrInvalidCity):
				errorJSONResponse(w, http.StatusBadRequest, "invalid city")
			case errors.As(err, &valErr):
				errorJSONResponse(w, http.StatusBadRequest, "validation error")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		jsonResponse(w, http.StatusCreated, pvz)
	}
}

func GetPVZListHandler(logger *zap.SugaredLogger, service pvzService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http_handlers.pvz.GetPVZListHandler"
		ctx := r.Context()
		logger := logger.With("req_id", middleware.GetReqID(ctx))

		var startDate, endDate *string
		if sd := r.URL.Query().Get("startDate"); sd != "" {
			startDate = &sd
		}

		if ed := r.URL.Query().Get("endDate"); ed != "" {
			endDate = &ed
		}

		page := 1
		if p := r.URL.Query().Get("page"); p != "" {
			if val, err := strconv.Atoi(p); err == nil && val > 0 {
				page = val
			}
		}

		limit := 10
		if l := r.URL.Query().Get("limit"); l != "" {
			if val, err := strconv.Atoi(l); err == nil && val > 0 && val <= 30 {
				limit = val
			}
		}

		pvzList, err := service.GetPVZList(ctx, startDate, endDate, page, limit)
		if err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusInternalServerError, "internal server error")
			return
		}

		jsonResponse(w, http.StatusOK, pvzList)
	}
}
