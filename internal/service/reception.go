package service

import (
	"assignment/internal/models"
	"context"
	"fmt"
)

func (s *Service) CreateReception(ctx context.Context, req models.NewReceptionRequest) (models.Reception, error) {
	const op = "internal.service.reception.CreateReception"

	if err := s.validator.Struct(req); err != nil {
		return models.Reception{}, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.repo.CreateReception(ctx, req)
	if err != nil {
		return models.Reception{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *Service) CloseLastReception(ctx context.Context, pvzID string) (models.Reception, error) {
	const op = "internal.service.reception.CloseLastReception"

	if pvzID == "" {
		return models.Reception{}, fmt.Errorf("%s: %w", op, ErrPVZIDRequired)
	}

	result, err := s.repo.CloseLastReception(ctx, pvzID)
	if err != nil {
		return models.Reception{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
