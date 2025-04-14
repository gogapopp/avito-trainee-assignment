package postgres

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"fmt"
)

func (s *Storage) CreatePVZ(ctx context.Context, pvz models.PVZ) (models.PVZ, error) {
	const op = "internal.repository.postgres.pvz.CreatePVZ"

	validCities := map[string]bool{
		"Москва":          true,
		"Санкт-Петербург": true,
		"Казань":          true,
	}

	if !validCities[pvz.City] {
		return models.PVZ{}, fmt.Errorf("%s: %w", op, repository.ErrInvalidCity)
	}

	pvzID := s.generateUUID()
	registrationDate := s.getCurrentTime()

	_, err := s.DB.Exec(ctx,
		"INSERT INTO pvz(id, registration_date, city) VALUES($1, $2, $3)",
		pvzID, registrationDate, pvz.City,
	)
	if err != nil {
		return models.PVZ{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.PVZ{
		ID:               pvzID,
		RegistrationDate: registrationDate,
		City:             pvz.City,
	}, nil
}

func (s *Storage) GetPVZList(ctx context.Context, startDate, endDate *string, page, limit int) ([]models.PVZWithReceptions, error) {
	const op = "internal.repository.postgres.pvz.GetPVZList"

	offset := (page - 1) * limit

	query := `
		SELECT p.id, p.registration_date, p.city
		FROM pvz p
	`

	args := []any{}
	if startDate != nil && endDate != nil {
		query += `
			WHERE EXISTS (
				SELECT 1 FROM reception r
				WHERE r.pvz_id = p.id
				AND r.date_time BETWEEN $1 AND $2
			)
		`
		args = append(args, *startDate, *endDate)
	}

	query += ` ORDER BY p.registration_date DESC LIMIT $` + fmt.Sprint(len(args)+1) + ` OFFSET $` + fmt.Sprint(len(args)+2)
	args = append(args, limit, offset)

	rows, err := s.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var pvzList []models.PVZWithReceptions

	for rows.Next() {
		var pvz models.PVZ
		if err := rows.Scan(&pvz.ID, &pvz.RegistrationDate, &pvz.City); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		receptions, err := s.getReceptionsByPVZID(ctx, pvz.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		pvzList = append(pvzList, models.PVZWithReceptions{
			PVZ:        pvz,
			Receptions: receptions,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows iteration: %w", op, err)
	}

	return pvzList, nil
}

func (s *Storage) getReceptionsByPVZID(ctx context.Context, pvzID string) ([]models.ReceptionProducts, error) {
	const op = "internal.repository.postgres.pvz.getReceptionsByPVZID"

	rows, err := s.DB.Query(ctx, `
		SELECT id, date_time, pvz_id, status
		FROM reception
		WHERE pvz_id = $1
		ORDER BY date_time DESC
	`, pvzID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var receptions []models.ReceptionProducts

	for rows.Next() {
		var reception models.Reception
		if err := rows.Scan(&reception.ID, &reception.DateTime, &reception.PVZID, &reception.Status); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		products, err := s.getProductsByReceptionID(ctx, reception.ID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		receptions = append(receptions, models.ReceptionProducts{
			Reception: reception,
			Products:  products,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return receptions, nil
}

func (s *Storage) getProductsByReceptionID(ctx context.Context, receptionID string) ([]models.Product, error) {
	const op = "internal.repository.postgres.pvz.getProductsByReceptionID"

	rows, err := s.DB.Query(ctx, `
		SELECT id, date_time, type, reception_id
		FROM product
		WHERE reception_id = $1
		ORDER BY date_time ASC
	`, receptionID)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.DateTime, &product.Type, &product.ReceptionID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return products, nil
}
