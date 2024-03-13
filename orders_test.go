package sddb

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"testing"
	"time"
)

type testStoragerOrder interface {
	OrderStorager

	clean(ctx context.Context) error
}

type config struct {
	ConnectTimeout time.Duration
	QueryTimeout   time.Duration
	Username       string
	Password       string
	DBName         string
	MigrationVer   int

	Host string
	Port int
}

type PostrgresTestSuiteOrder struct {
	suite.Suite
	testStoragerOrder

	tc  *tcpostgres.PostgresContainer
	cfg *config
}

func (ts *PostrgresTestSuiteOrder) SetupSuite() {
	cfg := &config{
		ConnectTimeout: 5 * time.Second,
		QueryTimeout:   5 * time.Second,
		Username:       "postgres",
		Password:       "test",
		DBName:         "postgres",
		MigrationVer:   1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pgc, err := tcpostgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:latest"),
		tcpostgres.WithDatabase(cfg.DBName),
		tcpostgres.WithUsername(cfg.Username),
		tcpostgres.WithPassword(cfg.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connection").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	require.NoError(ts.T(), err)

	cfg.Host, err = pgc.Host(ctx)
	require.NoError(ts.T(), err)

	port, err := pgc.MappedPort(ctx, "5432")
	require.NoError(ts.T(), err)

	cfg.Port = port.Int()

	ts.tc = pgc
	ts.cfg = cfg

	database_dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	log.Println(database_dsn)
	db, err := sqlx.Connect("postgres", database_dsn)

	storage := NewOrdersPostgresStorage(db)
	ts.testStoragerOrder = storage

	err = Migrate(db)
	require.NoError(ts.T(), err)

	ts.T().Logf("stared postgres at %s:%d", cfg.Host, cfg.Port)

}

// деструктор delete container
func (ts *PostrgresTestSuiteOrder) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoError(ts.T(), ts.tc.Terminate(ctx))
}

func testPostgres(t *testing.T) {
	suite.Run(t, new(PostrgresTestSuite))
}

func (ts *PostrgresTestSuiteOrder) testDummy() {}

func (s *OrdersPostgresStorage) clean(ctx context.Context) error {
	Newctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(Newctx, "DELETE FROM products")
	return err
}

func (ts *PostrgresTestSuiteOrder) SetupTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *PostrgresTestSuiteOrder) TearDownTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *PostrgresTestSuiteOrder) TestAddOrderGetOrder() {
	order := Orders{
		ID:            1,
		TgID:          1,
		UserName:      "11",
		FirstName:     "22",
		LastName:      "33",
		StatusOrder:   2,
		Pvz:           "33",
		Order:         "444",
		CreatedAt:     time.Now(),
		UpdateAt:      time.Now(),
		TypeDostavka:  0,
		PriceDelivery: 0,
		PriceFull:     0,
	}

	err := ts.AddOrder(context.Background(), order)
	ts.NoError(err)
	getOrders, err := ts.GetOrdersByTgID(context.Background(), 1)
	ts.NoError(err)
	if len(getOrders) >= 1 {
		ts.T().Error("len(getOrders) > 1")
	}
	ts.Require().Equal(order, getOrders[0], "getOrders[0] != order")

}
