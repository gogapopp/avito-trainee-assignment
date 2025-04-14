package service

import (
	"assignment/internal/models"
	"context"
	"fmt"
)

func (s *Service) AddProduct(ctx context.Context, req models.NewProductRequest) (models.Product, error) {
	const op = "internal.service.product.AddProduct"

	if err := s.validator.Struct(req); err != nil {
		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.repo.AddProduct(ctx, req)
	if err != nil {
		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *Service) DeleteLastProduct(ctx context.Context, pvzID string) error {
	const op = "internal.service.product.DeleteLastProduct"

	if pvzID == "" {
		return fmt.Errorf("%s: %w", op, ErrPVZIDRequired)
	}

	if err := s.repo.DeleteLastProduct(ctx, pvzID); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
