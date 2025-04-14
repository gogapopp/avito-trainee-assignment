package service

import (
	"assignment/internal/libs/jwt"
	"assignment/internal/models"
	"context"
	"fmt"
)

func (s *Service) RegisterUser(ctx context.Context, req models.AuthRequest) (models.User, error) {
	const op = "internal.service.auth.RegisterUser"

	if err := s.validator.Struct(req); err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := s.repo.RegisterUser(ctx, req)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return models.User{
		ID:    userID,
		Email: req.Email,
		Role:  req.Role,
	}, nil
}

func (s *Service) LoginUser(ctx context.Context, email, password string) (string, error) {
	const op = "internal.service.auth.LoginUser"

	userID, role, err := s.repo.LoginUser(ctx, email, password)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.GenerateToken(s.jwtSecret, userID, role)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (s *Service) DummyLogin(ctx context.Context, role string) (string, error) {
	const op = "internal.service.auth.DummyLogin"

	if role != "client" && role != "moderator" {
		return "", fmt.Errorf("%s: %w", op, ErrInvalidRole)
	}

	dummyUserID := "dummyid"
	token, err := jwt.GenerateToken(s.jwtSecret, dummyUserID, role)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}
