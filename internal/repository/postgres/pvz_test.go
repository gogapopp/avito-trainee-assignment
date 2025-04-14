package postgres_test

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPVZFunctions(t *testing.T) {
	ctx := context.Background()
	clearDatabase(t, ctx)

	t.Run("CreatePVZ", func(t *testing.T) {
		pvz := models.PVZ{
			City: "Москва",
		}

		createdPVZ, err := testStorage.CreatePVZ(ctx, pvz)
		require.NoError(t, err)
		assert.NotEmpty(t, createdPVZ.ID)
		assert.Equal(t, "Москва", createdPVZ.City)
		assert.False(t, createdPVZ.RegistrationDate.IsZero())

		invalidPVZ := models.PVZ{
			City: "Неизвестный город",
		}
		_, err = testStorage.CreatePVZ(ctx, invalidPVZ)
		assert.ErrorIs(t, err, repository.ErrInvalidCity)
	})

	t.Run("GetPVZList", func(t *testing.T) {
		cities := []string{"Санкт-Петербург", "Казань"}
		for _, city := range cities {
			pvz := models.PVZ{City: city}
			_, err := testStorage.CreatePVZ(ctx, pvz)
			require.NoError(t, err)
		}

		page, limit := 1, 10
		pvzList, err := testStorage.GetPVZList(ctx, nil, nil, page, limit)
		require.NoError(t, err)
		assert.Len(t, pvzList, 3)

		page, limit = 1, 2
		pvzList, err = testStorage.GetPVZList(ctx, nil, nil, page, limit)
		require.NoError(t, err)
		assert.Len(t, pvzList, 2)

		page, limit = 2, 2
		pvzList, err = testStorage.GetPVZList(ctx, nil, nil, page, limit)
		require.NoError(t, err)
		assert.Len(t, pvzList, 1)
	})
}
