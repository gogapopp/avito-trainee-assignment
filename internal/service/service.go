package service

import (
	"assignment/internal/models"
	"context"
	"errors"

	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

var (
	ErrInvalidRole   = errors.New("invalid role")
	ErrPVZIDRequired = errors.New("pvz id is required")
)

type Repository interface {
	RegisterUser(ctx context.Context, user models.AuthRequest) (string, error)
	LoginUser(ctx context.Context, email, password string) (string, string, error)

	CreatePVZ(ctx context.Context, pvz models.PVZ) (models.PVZ, error)
	GetPVZList(ctx context.Context, startDate, endDate *string, page, limit int) ([]models.PVZWithReceptions, error)

	CreateReception(ctx context.Context, reception models.NewReceptionRequest) (models.Reception, error)
	CloseLastReception(ctx context.Context, pvzID string) (models.Reception, error)

	AddProduct(ctx context.Context, product models.NewProductRequest) (models.Product, error)
	DeleteLastProduct(ctx context.Context, pvzID string) error
}

type Service struct {
	repo      Repository
	validator *validator.Validate
	jwtSecret string
	logger    *zap.SugaredLogger
}

func New(repo Repository, jwtSecret string, logger *zap.SugaredLogger) *Service {
	return &Service{
		repo:      repo,
		validator: validator.New(),
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}
