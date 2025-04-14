package postgres

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) AddProduct(ctx context.Context, req models.NewProductRequest) (models.Product, error) {
	const op = "internal.repository.postgres.product.AddProduct"

	var receptionID string
	err := s.DB.QueryRow(ctx, `
		SELECT id 
		FROM reception 
		WHERE pvz_id = $1 AND status = 'in_progress'
	`, req.PVZID).Scan(&receptionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, fmt.Errorf("%s: %w", op, repository.ErrNoActiveReception)
		}
		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	productID := s.generateUUID()
	now := s.getCurrentTime()

	_, err = s.DB.Exec(ctx, `
		INSERT INTO product(id, date_time, type, reception_id) 
		VALUES($1, $2, $3, $4)
	`, productID, now, req.Type, receptionID)
	if err != nil {
		return models.Product{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.Product{
		ID:          productID,
		DateTime:    now,
		Type:        req.Type,
		ReceptionID: receptionID,
	}, nil
}

func (s *Storage) DeleteLastProduct(ctx context.Context, pvzID string) error {
	const op = "internal.repository.postgres.product.DeleteLastProduct"

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	var receptionID string
	err = tx.QueryRow(ctx, `
		SELECT id 
		FROM reception 
		WHERE pvz_id = $1 AND status = 'in_progress'
	`, pvzID).Scan(&receptionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, repository.ErrNoActiveReception)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	var productID string
	err = tx.QueryRow(ctx, `
		SELECT id 
		FROM product 
		WHERE reception_id = $1 
		ORDER BY date_time DESC 
		LIMIT 1
	`, receptionID).Scan(&productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, repository.ErrNoProductsToDelete)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = tx.Exec(ctx, "DELETE FROM product WHERE id = $1", productID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
