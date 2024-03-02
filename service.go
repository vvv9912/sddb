package sddb

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type ProductsStorager interface {
	AddProduct(ctx context.Context, product Products) error

	ChangeProductByArticle(ctx context.Context, product Products) error

	GetCatalogNames(ctx context.Context) ([]string, error)

	GetAllProducts(ctx context.Context) ([]Products, error)

	GetProductsByCatalogIsAvailable(ctx context.Context, ctlg string) ([]Products, error)

	GetProductsByCatalog(ctx context.Context, ctlg string) ([]Products, error)

	GetProductByArticle(ctx context.Context, article int) (Products, error)
}
type StorageProduct struct {
	ProductsStorager
}

func NewStorageProducts(db *sqlx.DB) *StorageProduct {
	return &StorageProduct{ProductsStorager: NewProductsPostgresStorage(db)}
}

type CorzinaStorager interface {
	AddShopCart(ctx context.Context, shopCart ShopCart) error
	GetShopCartByTgID(ctx context.Context, tgId int64) ([]ShopCart, error)
	GetShopCartByTgIdAndArticle(ctx context.Context, tgId int64, article int) (ShopCart, error)
	UpdateShopCartByTgId(ctx context.Context, tgId int64, article int, quantity int) error
	GetShopCartDetailByTgId(ctx context.Context, tgId int64) ([]DbCorzineCatalog, error)
	DeleteShopCartByTgId(ctx context.Context, tgId int64) error
	DeleteShopCartByTgIdAndArticle(ctx context.Context, tgId int64, article int) error
}

type StorageCorzina struct {
	CorzinaStorager
}

func NewStorageCorzina(db *sqlx.DB) *StorageCorzina {
	return &StorageCorzina{CorzinaStorager: NewShopCartPostgresStorage(db)}
}

type OrderStorager interface {
	GetOrdersByTgID(ctx context.Context, tgId int64) ([]Orders, error)
	AddOrder(ctx context.Context, order Orders) error
}
type StorageOrder struct {
	OrderStorager
}

func NewStorageOrder(db *sqlx.DB) *StorageOrder {
	return &StorageOrder{OrderStorager: NewOrdersPostgresStorage(db)}
}

type StorageUser struct {
	UsersStorager
}

type UsersStorager interface {
	GetStatusUserByTgID(ctx context.Context, tgID int64) (int, int, error)
	AddUser(ctx context.Context, users Users) error
	UpdateStateByTgID(ctx context.Context, tgId int64, state int) error
	//GetCorzinaByTgID(ctx context.Context, tgID int64) ([]int64, error)
	//UpdateShopCartByTgId(ctx context.Context, tgId int64, corzina []int64) error
}

func NewStorageUser(db *sqlx.DB) *StorageUser {
	return &StorageUser{UsersStorager: NewUsersPostgresStorage(db)}

}
