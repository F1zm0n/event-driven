package repository_test

import (
	"context"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	pg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/F1zm0n/event-driven/auth/internal/entity"
	"github.com/F1zm0n/event-driven/auth/internal/repository"
)

type RepositoryTestSuite struct {
	suite.Suite
	DB    *sqlx.DB
	M     *migrate.Migrate
	cont  *pg.PostgresContainer
	repo  repository.CustomerRepository
	faker *gofakeit.Faker
	ctx   context.Context
}

const (
	DBName     = "auth-test"
	DBUser     = "postgres"
	DBPassword = "password"
)

func (r *RepositoryTestSuite) SetupTest() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background()
	postgresContainer, err := pg.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16-alpine"),
		pg.WithDatabase(DBName),
		pg.WithUsername(DBUser),
		pg.WithPassword(DBPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
			// WithStartupTimeout(10*time.Second)),
		))
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}
	r.cont = postgresContainer

	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(r.T(), err)

	r.DB, err = sqlx.Open("postgres", dsn)
	require.NoError(r.T(), err)

	_, filename, _, _ := runtime.Caller(0)
	migrationPath := "file://" + filepath.Join(filepath.Dir(filename), "../../migrations")
	r.T().Log(migrationPath)
	driver, err := postgres.WithInstance(r.DB.DB, &postgres.Config{DatabaseName: DBName})
	require.NoError(r.T(), err)

	r.M, err = migrate.NewWithDatabaseInstance(migrationPath, DBName, driver)
	require.NoError(r.T(), err)
	require.NoError(r.T(), r.M.Up())

	r.repo = repository.NewCustomerRepository(logger, r.DB)
	r.faker = gofakeit.New(0)

	r.ctx = context.Background()
}

func (r *RepositoryTestSuite) TearDownTest() {
	require.NoError(r.T(), r.M.Down())

	_, err := r.M.Close()
	require.NoError(r.T(), err)

	require.NoError(r.T(), r.cont.Terminate(context.Background()))
}

func (r *RepositoryTestSuite) TestCreateCustomer() {
	customer := entity.CustomerEntity{
		CustomerID: uuid.New(),
		Email:      r.faker.Email(),
		Password:   []byte(r.faker.Password(true, true, true, false, false, 8)),
	}
	err := r.repo.Create(r.ctx, customer)
	require.NoError(r.T(), err)
}

func (r *RepositoryTestSuite) TestGetCustomerByID() {
	customer := entity.CustomerEntity{
		CustomerID: uuid.New(),
		Email:      r.faker.Email(),
		Password:   []byte(r.faker.Password(true, true, true, false, false, 8)),
	}
	err := r.repo.Create(r.ctx, customer)
	require.NoError(r.T(), err)

	getCustomer, err := r.repo.GetByID(r.ctx, customer.CustomerID)
	require.NoError(r.T(), err)
	require.Equal(r.T(), getCustomer, customer)
}

func (r *RepositoryTestSuite) TestGetCustomerByEmail() {
	customer := entity.CustomerEntity{
		CustomerID: uuid.New(),
		Email:      r.faker.Email(),
		Password:   []byte(r.faker.Password(true, true, true, false, false, 8)),
	}
	err := r.repo.Create(r.ctx, customer)
	require.NoError(r.T(), err)

	getCustomer, err := r.repo.GetByEmail(r.ctx, customer.Email)
	require.NoError(r.T(), err)
	require.Equal(r.T(), getCustomer, customer)
}

func TestSomethingRuns(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}
