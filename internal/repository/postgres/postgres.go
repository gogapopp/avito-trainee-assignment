package postgres

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	DB         *pgxpool.Pool
	passSecret string
}

func New(ctx context.Context, dsn, passSecret string) (*Storage, error) {
	const op = "internal.repository.postgres.postgres.New"

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &Storage{
		DB:         pool,
		passSecret: passSecret,
	}, nil
}

func (s *Storage) generatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(s.passSecret)))
}

func (s *Storage) generateUUID() string {
	return uuid.New().String()
}

func (s *Storage) getCurrentTime() time.Time {
	return time.Now().UTC()
}
