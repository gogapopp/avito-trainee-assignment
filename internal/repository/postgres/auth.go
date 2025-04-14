package postgres

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) RegisterUser(ctx context.Context, user models.AuthRequest) (string, error) {
	const op = "internal.repository.postgres.auth.RegisterUser"

	userID := s.generateUUID()
	hashedPassword := s.generatePasswordHash(user.Password)

	_, err := s.DB.Exec(ctx,
		"INSERT INTO users(id, email, password_hash, role) VALUES($1, $2, $3, $4)",
		userID, user.Email, hashedPassword, user.Role,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return "", fmt.Errorf("%s: %w", op, repository.ErrUserExists)
			}
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (s *Storage) LoginUser(ctx context.Context, email, password string) (string, string, error) {
	const op = "internal.repository.postgres.auth.LoginUser"

	hashedPassword := s.generatePasswordHash(password)

	var (
		userID       string
		passwordHash string
		role         string
	)

	err := s.DB.QueryRow(ctx,
		"SELECT id, password_hash, role FROM users WHERE email = $1",
		email,
	).Scan(&userID, &passwordHash, &role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", fmt.Errorf("%s: %w", op, repository.ErrInvalidCredentials)
		}
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	if passwordHash != hashedPassword {
		return "", "", fmt.Errorf("%s: %w", op, repository.ErrInvalidCredentials)
	}

	return userID, role, nil
}
