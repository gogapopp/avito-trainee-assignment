package service

import (
	"assignment/internal/models"
	"context"
	"fmt"
)

func (s *Service) CreatePVZ(ctx context.Context, pvz models.PVZ) (models.PVZ, error) {
	const op = "internal.service.pvz.CreatePVZ"

	if err := s.validator.Struct(pvz); err != nil {
		return models.PVZ{}, fmt.Errorf("%s: %w", op, err)
	}

	result, err := s.repo.CreatePVZ(ctx, pvz)
	if err != nil {
		return models.PVZ{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (s *Service) GetPVZList(ctx context.Context, startDate, endDate *string, page, limit int) ([]models.PVZWithReceptions, error) {
	const op = "internal.service.pvz.GetPVZList"
	// set default value
	if page <= 0 {
		page = 1
	}
	// set default value
	if limit <= 0 || limit > 30 {
		limit = 10
	}

	result, err := s.repo.GetPVZList(ctx, startDate, endDate, page, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
