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

type productService interface {
	AddProduct(ctx context.Context, req models.NewProductRequest) (models.Product, error)
	DeleteLastProduct(ctx context.Context, pvzID string) error
}

func AddProductHandler(logger *zap.SugaredLogger, service productService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http_handlers.product.AddProductHandler"
		ctx := r.Context()
		logger := logger.With("req_id", middleware.GetReqID(ctx))

		var req models.NewProductRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Errorf("%s: %w", op, err)
			errorJSONResponse(w, http.StatusBadRequest, "invalid request body")
			return
		}

		product, err := service.AddProduct(ctx, req)
		if err != nil {
			var valErr validator.ValidationErrors
			logger.Errorf("%s: %w", op, err)

			switch {
			case errors.Is(err, repository.ErrNoActiveReception):
				errorJSONResponse(w, http.StatusBadRequest, "no active reception found")
			case errors.As(err, &valErr):
				errorJSONResponse(w, http.StatusBadRequest, "validation error")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		jsonResponse(w, http.StatusCreated, product)
	}
}

func DeleteLastProductHandler(logger *zap.SugaredLogger, service productService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.server.http_handlers.product.DeleteLastProductHandler"
		ctx := r.Context()
		logger := logger.With("req_id", middleware.GetReqID(ctx))

		pvzID := chi.URLParam(r, "pvzId")
		if pvzID == "" {
			errorJSONResponse(w, http.StatusBadRequest, "pvz id is required")
			return
		}

		err := service.DeleteLastProduct(ctx, pvzID)
		if err != nil {
			logger.Errorf("%s: %w", op, err)

			switch {
			case errors.Is(err, repository.ErrNoActiveReception):
				errorJSONResponse(w, http.StatusBadRequest, "no active reception found")
			case errors.Is(err, repository.ErrNoProductsToDelete):
				errorJSONResponse(w, http.StatusBadRequest, "no products to delete")
			default:
				errorJSONResponse(w, http.StatusInternalServerError, "internal server error")
			}
			return
		}

		jsonResponse(w, http.StatusOK, nil)
	}
}
