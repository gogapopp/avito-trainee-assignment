package postgres_test

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthFunctions(t *testing.T) {
	ctx := context.Background()
	clearDatabase(t, ctx)

	t.Run("RegisterUser", func(t *testing.T) {
		user := models.AuthRequest{
			Email:    "test@example.com",
			Password: "password123",
			Role:     "moderator",
		}

		userID, err := testStorage.RegisterUser(ctx, user)
		require.NoError(t, err)
		assert.NotEmpty(t, userID)

		_, err = testStorage.RegisterUser(ctx, user)
		assert.ErrorIs(t, err, repository.ErrUserExists)
	})

	t.Run("LoginUser", func(t *testing.T) {
		// successful login
		userID, role, err := testStorage.LoginUser(ctx, "test@example.com", "password123")
		require.NoError(t, err)
		assert.NotEmpty(t, userID)
		assert.Equal(t, "client", role)

		// wrong password
		_, _, err = testStorage.LoginUser(ctx, "test@example.com", "wrongpassword")
		assert.ErrorIs(t, err, repository.ErrInvalidCredentials)

		// wrong email
		_, _, err = testStorage.LoginUser(ctx, "nonexistent@example.com", "password123")
		assert.ErrorIs(t, err, repository.ErrInvalidCredentials)
	})
}
