package postgres_test

import (
	"assignment/internal/models"
	"assignment/internal/repository"
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductFunctions(t *testing.T) {
	ctx := context.Background()
	clearDatabase(t, ctx)

	pvz := models.PVZ{City: "Москва"}
	createdPVZ, err := testStorage.CreatePVZ(ctx, pvz)
	require.NoError(t, err)
	pvzID := createdPVZ.ID

	req := models.NewReceptionRequest{PVZID: pvzID}
	reception, err := testStorage.CreateReception(ctx, req)
	require.NoError(t, err)

	t.Run("AddProduct", func(t *testing.T) {
		productReq := models.NewProductRequest{
			PVZID: pvzID,
			Type:  "TestProduct",
		}

		product, err := testStorage.AddProduct(ctx, productReq)
		require.NoError(t, err)
		assert.NotEmpty(t, product.ID)
		assert.Equal(t, "TestProduct", product.Type)
		assert.Equal(t, reception.ID, product.ReceptionID)

		for i := 0; i < 3; i++ {
			productReq := models.NewProductRequest{
				PVZID: pvzID,
				Type:  fmt.Sprintf("TestProduct%d", i),
			}
			_, err := testStorage.AddProduct(ctx, productReq)
			require.NoError(t, err)
		}

		invalidReq := models.NewProductRequest{
			PVZID: uuid.New().String(),
			Type:  "InvalidProduct",
		}
		_, err = testStorage.AddProduct(ctx, invalidReq)
		assert.ErrorIs(t, err, repository.ErrNoActiveReception)
	})

	t.Run("DeleteLastProduct", func(t *testing.T) {
		err := testStorage.DeleteLastProduct(ctx, pvzID)
		require.NoError(t, err)

		err = testStorage.DeleteLastProduct(ctx, pvzID)
		require.NoError(t, err)
		err = testStorage.DeleteLastProduct(ctx, pvzID)
		require.NoError(t, err)
		err = testStorage.DeleteLastProduct(ctx, pvzID)
		require.NoError(t, err)

		err = testStorage.DeleteLastProduct(ctx, pvzID)
		assert.ErrorIs(t, err, repository.ErrNoProductsToDelete)

		_, err = testStorage.CloseLastReception(ctx, pvzID)
		require.NoError(t, err)

		err = testStorage.DeleteLastProduct(ctx, pvzID)
		assert.ErrorIs(t, err, repository.ErrNoActiveReception)
	})
}
