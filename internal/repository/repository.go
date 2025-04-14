package repository

import "errors"

var (
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrUserExists           = errors.New("user already exists")
	ErrInvalidCity          = errors.New("invalid city")
	ErrNoActiveReception    = errors.New("no active reception")
	ErrReceptionClosed      = errors.New("reception already closed")
	ErrActiveReceptionExist = errors.New("active reception already exists")
	ErrNoProductsToDelete   = errors.New("no products to delete")
	ErrPVZNotFound          = errors.New("pvz not found")
)
