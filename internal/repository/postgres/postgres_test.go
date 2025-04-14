package postgres_test

import (
	"assignment/internal/models"
	"assignment/internal/repository/postgres"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testStorage *postgres.Storage
	pgDSN       string
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15-alpine",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=testdb",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			fmt.Println(err)
		}
	}()

	hostAndPort := resource.GetHostPort("5432/tcp")
	pgDSN = fmt.Sprintf("postgres://postgres:postgres@%s/testdb?sslmode=disable", hostAndPort)

	pool.MaxWait = 60 * time.Second
	if err = pool.Retry(func() error {
		ctx := context.Background()
		db, err := pgxpool.New(ctx, pgDSN)
		if err != nil {
			return err
		}
		return db.Ping(ctx)
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx := context.Background()
	testStorage, err = postgres.New(ctx, pgDSN, "test-secret")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = setupTestSchema(ctx, testStorage.DB)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	code := m.Run()

	testStorage.DB.Close()
	os.Exit(code)
}

func setupTestSchema(ctx context.Context, db *pgxpool.Pool) error {
	_, err := db.Exec(ctx, `
		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL
		);
		
		CREATE TABLE pvz (
			id TEXT PRIMARY KEY,
			registration_date TIMESTAMP WITH TIME ZONE NOT NULL,
			city TEXT NOT NULL
		);
		
		CREATE TABLE reception (
			id TEXT PRIMARY KEY,
			date_time TIMESTAMP WITH TIME ZONE NOT NULL,
			pvz_id TEXT NOT NULL REFERENCES pvz(id) ON DELETE CASCADE,
			status TEXT NOT NULL CHECK (status IN ('in_progress', 'close'))
		);
		
		CREATE TABLE product (
			id TEXT PRIMARY KEY,
			date_time TIMESTAMP WITH TIME ZONE NOT NULL,
			type TEXT NOT NULL,
			reception_id TEXT NOT NULL REFERENCES reception(id) ON DELETE CASCADE
		);
		
		CREATE INDEX reception_pvz_id_idx ON reception(pvz_id);
		CREATE INDEX reception_status_idx ON reception(status);
		CREATE INDEX product_reception_id_idx ON product(reception_id);
	`)

	return err
}

func clearDatabase(t *testing.T, ctx context.Context) {
	tables := []string{"product", "reception", "pvz", "users"}
	for _, table := range tables {
		_, err := testStorage.DB.Exec(ctx, fmt.Sprintf("DELETE FROM %s", table))
		require.NoError(t, err)
	}
}

func TestFullWorkflow(t *testing.T) {
	ctx := context.Background()
	clearDatabase(t, ctx)

	pvz := models.PVZ{City: "Москва"}
	createdPVZ, err := testStorage.CreatePVZ(ctx, pvz)
	require.NoError(t, err)
	pvzID := createdPVZ.ID

	receptionReq := models.NewReceptionRequest{PVZID: pvzID}
	reception, err := testStorage.CreateReception(ctx, receptionReq)
	require.NoError(t, err)
	assert.Equal(t, "in_progress", reception.Status)

	productTypes := []string{"Phone", "Laptop", "Tablet"}
	for _, productType := range productTypes {
		productReq := models.NewProductRequest{
			PVZID: pvzID,
			Type:  productType,
		}
		product, err := testStorage.AddProduct(ctx, productReq)
		require.NoError(t, err)
		assert.Equal(t, productType, product.Type)
	}

	err = testStorage.DeleteLastProduct(ctx, pvzID)
	require.NoError(t, err)

	closedReception, err := testStorage.CloseLastReception(ctx, pvzID)
	require.NoError(t, err)
	assert.Equal(t, "close", closedReception.Status)

	pvzList, err := testStorage.GetPVZList(ctx, nil, nil, 1, 10)
	require.NoError(t, err)
	assert.Len(t, pvzList, 1)

	assert.Len(t, pvzList[0].Receptions, 1)
	assert.Equal(t, "close", pvzList[0].Receptions[0].Reception.Status)

	assert.Len(t, pvzList[0].Receptions[0].Products, 2)
	assert.Equal(t, "Phone", pvzList[0].Receptions[0].Products[0].Type)
	assert.Equal(t, "Laptop", pvzList[0].Receptions[0].Products[1].Type)
}
