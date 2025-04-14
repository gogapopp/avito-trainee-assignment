package postgres

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) CreateReception(ctx context.Context, req models.NewReceptionRequest) (models.Reception, error) {
	const op = "internal.repository.postgres.reception.CreateReception"

	var exists bool
	err := s.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM pvz WHERE id = $1)", req.PVZID).Scan(&exists)
	if err != nil {
		return models.Reception{}, fmt.Errorf("%s: check pvz: %w", op, err)
	}

	if !exists {
		return models.Reception{}, fmt.Errorf("%s: %w", op, repository.ErrPVZNotFound)
	}

	var activeReceptionExists bool
	err = s.DB.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM reception 
			WHERE pvz_id = $1 AND status = 'in_progress'
		)
	`, req.PVZID).Scan(&activeReceptionExists)
	if err != nil {
		return models.Reception{}, fmt.Errorf("%s: %w", op, err)
	}

	if activeReceptionExists {
		return models.Reception{}, fmt.Errorf("%s: %w", op, repository.ErrActiveReceptionExist)
	}

	receptionID := s.generateUUID()
	now := s.getCurrentTime()

	_, err = s.DB.Exec(ctx, `
		INSERT INTO reception(id, date_time, pvz_id, status) 
		VALUES($1, $2, $3, $4)
	`, receptionID, now, req.PVZID, "in_progress")
	if err != nil {
		return models.Reception{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.Reception{
		ID:       receptionID,
		DateTime: now,
		PVZID:    req.PVZID,
		Status:   "in_progress",
	}, nil
}

func (s *Storage) CloseLastReception(ctx context.Context, pvzID string) (models.Reception, error) {
	const op = "internal.repository.postgres.reception.CloseLastReception"

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return models.Reception{}, fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var reception models.Reception
	err = tx.QueryRow(ctx, `
		SELECT id, date_time, pvz_id, status 
		FROM reception 
		WHERE pvz_id = $1 AND status = 'in_progress'
	`, pvzID).Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Reception{}, fmt.Errorf("%s: %w", op, repository.ErrNoActiveReception)
		}
		return models.Reception{}, fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE reception 
		SET status = 'close' 
		WHERE id = $1
	`, reception.ID)
	if err != nil {
		return models.Reception{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return models.Reception{}, fmt.Errorf("%s: %w", op, err)
	}

	reception.Status = "close"

	return reception, nil
}
