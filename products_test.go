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

type testStorager interface {
	ProductsStorager

	clean(ctx context.Context) error
}

type Config struct {
	ConnectTimeout time.Duration
	QueryTimeout   time.Duration
	Username       string
	Password       string
	DBName         string
	MigrationVer   int

	Host string
	Port int
}

type PostrgresTestSuite struct {
	suite.Suite
	testStorager

	tc  *tcpostgres.PostgresContainer
	cfg *Config
}

func (ts *PostrgresTestSuite) SetupSuite() {
	cfg := &Config{
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

	storage := NewProductsPostgresStorage(db)

	ts.testStorager = storage

	err = Migrate(db)
	require.NoError(ts.T(), err)

	ts.T().Logf("stared postgres at %s:%d", cfg.Host, cfg.Port)

}

// деструктор delete container
func (ts *PostrgresTestSuite) TearDownSuite() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	require.NoError(ts.T(), ts.tc.Terminate(ctx))
}

func TestPostgres(t *testing.T) {
	suite.Run(t, new(PostrgresTestSuite))
}

func (ts *PostrgresTestSuite) TestDummy() {}

func (s *ProductsPostgresStorage) clean(ctx context.Context) error {
	Newctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(Newctx, "DELETE FROM products")
	return err
}

func (ts *PostrgresTestSuite) SetupTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *PostrgresTestSuite) TearDownTest() {
	ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *PostrgresTestSuite) TestAddProduct() {

	products := Products{
		Article:      1,
		Catalog:      "Test Catalog",
		Name:         "Name",
		Description:  "Descrip",
		PhotoUrl:     [][]byte{{1}, {2}, {3}},
		Price:        144,
		Length:       1,
		Width:        1,
		Height:       1,
		Weight:       1,
		Availability: true,
	}
	err := ts.AddProduct(context.Background(), products)
	ts.NoError(err)
	goods, err := ts.GetAllProducts(context.Background())
	// ts.GetCatalogNames(context.Background())
	ts.NoError(err)

	if len(goods) > 1 {
		ts.T().Error("len(goods) > 1")
	}
	ts.Require().Equal(products, goods[0], "goods[0] != products")
	//	log.Println(goods)
	//ts.Require().NoError(ts.clean(context.Background()))
}
func (ts *PostrgresTestSuite) TestGetProductsByCatalogIsAvailable() {

	products1 := Products{
		Article:      2,
		Catalog:      "Test Catalog",
		Name:         "Name",
		Description:  "Descrip",
		PhotoUrl:     [][]byte{{1}, {2}, {3}},
		Price:        144,
		Length:       1,
		Width:        1,
		Height:       1,
		Weight:       1,
		Availability: true,
	}
	products2 := Products{
		Article:      1,
		Catalog:      "Test Catalog",
		Name:         "Name",
		Description:  "Descrip",
		PhotoUrl:     [][]byte{{1}, {2}, {3}},
		Price:        144,
		Length:       1,
		Width:        1,
		Height:       1,
		Weight:       1,
		Availability: false,
	}
	err := ts.AddProduct(context.Background(), products1)
	ts.NoError(err)
	err = ts.AddProduct(context.Background(), products2)
	ts.NoError(err)
	goods, err := ts.GetProductsByCatalogIsAvailable(context.Background(), "Test Catalog")
	ts.NoError(err)

	if len(goods) > 1 {
		ts.T().Error("len(goods) > 1")
	}
	ts.Require().Equal(products1, goods[0], "goods[0] != products")
	//	log.Println(goods)
	//ts.Require().NoError(ts.clean(context.Background()))
}
func (ts *PostrgresTestSuite) TestGetProductByArticle() {

	products1 := Products{
		Article:      2,
		Catalog:      "Test Catalog",
		Name:         "Name",
		Description:  "Descrip",
		PhotoUrl:     [][]byte{{1}, {2}, {3}},
		Price:        144,
		Length:       1,
		Width:        1,
		Height:       1,
		Weight:       1,
		Availability: true,
	}
	products2 := Products{
		Article:      1,
		Catalog:      "Test Catalog",
		Name:         "Name",
		Description:  "Descrip",
		PhotoUrl:     [][]byte{{1}, {2}, {3}},
		Price:        144,
		Length:       1,
		Width:        1,
		Height:       1,
		Weight:       1,
		Availability: false,
	}
	err := ts.AddProduct(context.Background(), products1)
	ts.NoError(err)
	err = ts.AddProduct(context.Background(), products2)
	ts.NoError(err)
	goods, err := ts.GetProductByArticle(context.Background(), 1)
	ts.NoError(err)

	ts.Require().Equal(products2, goods, "goods != products")
	//	log.Println(goods)
	//ts.Require().NoError(ts.clean(context.Background()))
}

func (ts *PostrgresTestSuite) TestChangeProductByArticle() {

	products1 := Products{
		Article:      1,
		Catalog:      "Test Catalog11",
		Name:         "Name11",
		Description:  "Descrip",
		PhotoUrl:     [][]byte{{1}, {2}, {3}},
		Price:        144,
		Length:       1,
		Width:        1,
		Height:       1,
		Weight:       1,
		Availability: true,
	}
	products2 := Products{
		Article:      1,
		Catalog:      "Test Catalog22",
		Name:         "Name22",
		Description:  "Descrip22",
		PhotoUrl:     [][]byte{{1}, {2}, {3}},
		Price:        122,
		Length:       2,
		Width:        2,
		Height:       2,
		Weight:       2,
		Availability: false,
	}
	err := ts.AddProduct(context.Background(), products1)
	ts.NoError(err)

	err = ts.ChangeProductByArticle(context.Background(), products2)
	ts.NoError(err)

	goods, err := ts.GetProductByArticle(context.Background(), 1)
	ts.NoError(err)

	ts.Require().Equal(products2, goods, "goods != products")
	//	log.Println(goods)
	//ts.Require().NoError(ts.clean(context.Background()))
}
