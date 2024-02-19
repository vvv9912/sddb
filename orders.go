package main

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
)

type OrdersPostgresStorage struct {
	db *sqlx.DB
}

func NewOrdersPostgresStorage(db *sqlx.DB) *OrdersPostgresStorage {
	return &OrdersPostgresStorage{db: db}
}
func (s *OrdersPostgresStorage) OrdersByTgID(ctx context.Context, tgId int64) ([]Orders, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	var corzine []dbOrder
	if err := conn.SelectContext(ctx,
		&corzine,
		`SELECT
     			id AS o_id,
     			tg_id AS o_tg_id,
     			status_order AS o_status_order,
     			pvz AS o_pvz,
     			type_dostavka AS o_type_dostavka,
     			orderr AS o_order,
     			created_at AS o_created_at,
     			read_at AS o_read_at
	 			FROM orders
	 			WHERE tg_id = $1`,
		tgId); err != nil {
		return nil, err
	}
	//return lo.Map(corzine, func(corzin dbCorzine, _ int) model.Corzine { return model.Corzine(corzin) }), nil
	return lo.Map(corzine, func(corzin dbOrder, _ int) Orders { return Orders(corzin) }), nil
}

type dbOrder struct {
	ID            int       `db:"o_id"`
	TgID          int64     `db:"o_tg_id"`
	StatusOrder   int       `db:"o_status_order"`
	Pvz           string    `db:"o_pvz"`
	Order         string    `db:"o_order"`
	CreatedAt     time.Time `db:"o_created_at"`
	ReadAt        time.Time `db:"o_read_at"`
	TypeDostavka  int       `db:"o_type_dostavka"`
	PriceDelivery float64   `db:"o_price_delivery"`
	PriceFull     float64   `db:"o_price_full"`
}

func (s *OrdersPostgresStorage) OrdersByArticle() {

}

// todo
func (s *OrdersPostgresStorage) AddOrders(ctx context.Context, order Orders) error {
	//conn, err := s.db.Connx(ctx)
	// &sql.TxOptions{Isolation: sql.LevelSerializable}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO orders (tg_id, status_order,  orderr, created_at, read_at, pvz, type_dostavka)
	    				VALUES ($1, $2, $3, $4,$5,$6,$7)
	    				ON CONFLICT DO NOTHING;`,
		order.TgID,
		order.StatusOrder,
		//
		order.Order,
		order.CreatedAt,
		order.ReadAt,
		order.Pvz,
		order.TypeDostavka,
	); err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
