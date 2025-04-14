package postgres_test

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReceptionFunctions(t *testing.T) {
	ctx := context.Background()
	clearDatabase(t, ctx)

	pvz := models.PVZ{City: "Москва"}
	createdPVZ, err := testStorage.CreatePVZ(ctx, pvz)
	require.NoError(t, err)
	pvzID := createdPVZ.ID

	t.Run("CreateReception", func(t *testing.T) {
		req := models.NewReceptionRequest{
			PVZID: pvzID,
		}

		reception, err := testStorage.CreateReception(ctx, req)
		require.NoError(t, err)
		assert.NotEmpty(t, reception.ID)
		assert.Equal(t, pvzID, reception.PVZID)
		assert.Equal(t, "in_progress", reception.Status)

		_, err = testStorage.CreateReception(ctx, req)
		assert.ErrorIs(t, err, repository.ErrActiveReceptionExist)

		invalidReq := models.NewReceptionRequest{
			PVZID: uuid.New().String(),
		}
		_, err = testStorage.CreateReception(ctx, invalidReq)
		assert.ErrorIs(t, err, repository.ErrPVZNotFound)
	})

	t.Run("CloseLastReception", func(t *testing.T) {
		reception, err := testStorage.CloseLastReception(ctx, pvzID)
		require.NoError(t, err)
		assert.Equal(t, "close", reception.Status)

		_, err = testStorage.CloseLastReception(ctx, pvzID)
		assert.ErrorIs(t, err, repository.ErrNoActiveReception)
	})
}
