package rabbitmq_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/F1zm0n/event-driven/auth/internal/dto"
	"github.com/F1zm0n/event-driven/auth/internal/repository"
	"github.com/F1zm0n/event-driven/auth/internal/transport/rabbitmq"
)

type CustomerTest struct {
	suite.Suite
	q    amqp.Queue
	ch   *amqp.Channel
	conn *amqp.Connection
	repo repository.CustomerRepository
	db   *sqlx.DB
	ctx  context.Context
}

func (t *CustomerTest) SetupTest() {
	err := godotenv.Load("/home/F1zm0/Files/projects/gov1/event-driven/auth/.env")
	require.NoError(t.T(), err)
	t.conn, err = amqp.Dial(os.Getenv("RABBIT_DSN"))
	require.NoError(t.T(), err)

	t.ch, err = t.conn.Channel()
	require.NoError(t.T(), err)

	t.q = rabbitmq.MustMakeQueue(t.ch)

	var (
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

		db   = sqlx.MustOpen("postgres", os.Getenv("DSN"))
		repo = repository.NewCustomerRepository(logger, db)
	)
	t.repo = repo
	t.db = db
	rabbitmq.MustMakeExchangeBindings(t.ch, t.q)
}

func (t *CustomerTest) TearDownTest() {
	t.db.Close()
	t.ch.Close()
	t.conn.Close()
}

func (t *CustomerTest) TestRabbitMQEndpoint() {
	req := require.New(t.T())
	customerDto := dto.CustomerDto{
		Email:    "bimbam@gmail.com",
		Password: "nikker",
	}
	b, err := json.Marshal(customerDto)
	req.NoError(err)

	err = t.ch.PublishWithContext(
		context.Background(),
		os.Getenv("EXCHANGE_NAME"),
		t.q.Name,
		false,
		false,
		amqp.Publishing{
			Headers: amqp.Table{
				"operation": os.Getenv("REGISTER_HEADER"),
			},
			Body: b,
		},
	)
	req.NoError(err)
	time.Sleep(time.Second)
	customer, err := t.repo.GetByEmail(t.ctx, customerDto.Email)
	req.NoError(err)
	assert.Equal(t.T(), customerDto.Email, customer.Email)
}

func TestLaunchTests(t *testing.T) {
	suite.Run(t, new(CustomerTest))
}
