package sddb

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type ProductsStorager interface {
	AddProduct(ctx context.Context, product Products) error
	ChangeProductByArticle(ctx context.Context, product Products) error
	Catalog(ctx context.Context) ([]string, error)

	SelectAllProducts(ctx context.Context) ([]Products, error)

	ProductsByCatalog(ctx context.Context, ctlg string) ([]Products, error)

	ProductByArticle(ctx context.Context, article int) (Products, error)
}
type StorageProduct struct {
	ProductsStorager
}

func NewStorageProducts(db *sqlx.DB) *StorageProduct {
	return &StorageProduct{ProductsStorager: NewProductsPostgresStorage(db)}
}

type CorzinaStorager interface {
	AddCorzinas(ctx context.Context, corz Corzine) error
	CorzinaByTgId(ctx context.Context, tgId int64) ([]Corzine, error)
	CorzinaByTgIdANDAtricle(ctx context.Context, tgId int64, article int) (Corzine, error)
	UpdateCorzinaByTgId(ctx context.Context, tgId int64, article int, quantity int) error
	CorzinaByTgIdwithCalalog(ctx context.Context, tgId int64) ([]DbCorzineCatalog, error)
	DeleteCorzinaByTgID(ctx context.Context, tgId int64) error
	DeleteCorzinaByTgIDandArticle(ctx context.Context, tgId int64, article int) error
}

type StorageCorzina struct {
	CorzinaStorager
}

func NewStorageCorzina(db *sqlx.DB) *StorageCorzina {
	return &StorageCorzina{CorzinaStorager: NewCorzinaPostgresStorage(db)}
}

type OrderStorager interface {
	OrdersByTgID(ctx context.Context, tgId int64) ([]Orders, error)
	AddOrders(ctx context.Context, order Orders) error
}
type StorageOrder struct {
	OrderStorager
}

func NewStorageOrder(db *sqlx.DB) *StorageOrder {
	return &StorageOrder{OrderStorager: NewOrdersPostgresStorage(db)}
}
