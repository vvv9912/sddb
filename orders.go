package sddb

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
)

const (
	StatusOrderNew        = 1 //заказ новый
	StatusOrderRead       = 2 //заказ принят
	StatusOrderProcessing = 3 //заказ обрабатывается
	StatusOrderProcessed  = 4 //заказ обработан
	StatusOrderDelivered  = 5 //заказ доставляется
	StatusOrderCanceled   = 6 //заказ отменен
	StatusOrderReturned   = 7 //заказ возвращен
	StatusOrderComplete   = 8 //заказ выполнен
)

type OrdersPostgresStorage struct {
	db *sqlx.DB
}

func NewOrdersPostgresStorage(db *sqlx.DB) *OrdersPostgresStorage {
	return &OrdersPostgresStorage{db: db}
}

// use tg
func (s *OrdersPostgresStorage) GetOrdersByTgID(ctx context.Context, tgId int64) ([]Orders, error) {
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
	 			user_name AS o_user_name,
	 			first_name AS o_first_name,
	 			last_name AS o_last_name,
     			status_order AS o_status_order,
     			pvz AS o_pvz,
     			type_dostavka AS o_type_dostavka,
     			orderr AS o_order,
     			created_at  AT TIME ZONE 'UTC' AS o_created_at,
     			update_at  AT TIME ZONE 'UTC' AS o_update_at
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
	UserName      string    `db:"o_user_name"`
	FirstName     string    `db:"o_first_name"`
	LastName      string    `db:"o_last_name"`
	StatusOrder   int       `db:"o_status_order"`
	Pvz           string    `db:"o_pvz"`
	Order         string    `db:"o_order"`
	CreatedAt     time.Time `db:"o_created_at"`
	UpdateAt      time.Time `db:"o_update_at"`
	TypeDostavka  int       `db:"o_type_dostavka"`
	PriceDelivery float64   `db:"o_price_delivery"`
	PriceFull     float64   `db:"o_price_full"`
}

// use tg
func (s *OrdersPostgresStorage) AddOrder(ctx context.Context, order Orders) error {
	//conn, err := s.db.Connx(ctx)
	// &sql.TxOptions{Isolation: sql.LevelSerializable}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO orders (tg_id, 	user_name,first_name,last_name,status_order,  orderr, created_at, update_at, pvz, type_dostavka)
	    				VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9,$10)
	    				ON CONFLICT DO NOTHING;`,
		order.TgID,
		order.UserName,
		order.FirstName,
		order.LastName,
		order.StatusOrder,
		//
		order.Order,
		order.CreatedAt,
		order.UpdateAt,
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
func (s *OrdersPostgresStorage) GetOrderByStatus(ctx context.Context, statusOrder int) ([]Orders, error) {
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
     			user_name AS o_user_name,
     			first_name AS o_first_name,
	 			last_name AS o_last_name,
     			status_order AS o_status_order,
     			pvz AS o_pvz,
     			type_dostavka AS o_type_dostavka,
     			orderr AS o_order,
     			created_at  AT TIME ZONE 'UTC' AS o_created_at,
     			update_at  AT TIME ZONE 'UTC' AS o_update_at
	 			FROM orders
	 			WHERE status_order = $1`,
		statusOrder); err != nil {
		return nil, err
	}

	return lo.Map(corzine, func(corzin dbOrder, _ int) Orders { return Orders(corzin) }), nil
}

func (s *OrdersPostgresStorage) GetOrderByTimeAndStatus(ctx context.Context, statusOrder int, time2 time.Time) ([]Orders, error) {
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
     			user_name AS o_user_name,
     			first_name AS o_first_name,
	 			last_name AS o_last_name,
     			status_order AS o_status_order,
     			pvz AS o_pvz,
     			type_dostavka AS o_type_dostavka,
     			orderr AS o_order,
     			created_at  AT TIME ZONE 'UTC' AS o_created_at,
     				update_at  AT TIME ZONE 'UTC' AS o_update_at
	 			FROM orders
	 		    WHERE status_order = $1 AND created_at > $2`,
		statusOrder, time2); err != nil {
		return nil, err
	}

	return lo.Map(corzine, func(corzin dbOrder, _ int) Orders { return Orders(corzin) }), nil
}
func (s *OrdersPostgresStorage) UpdateOrderByStatus(ctx context.Context, statusOrder int, orderID int) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.ExecContext(ctx,
		`UPDATE orders
	 SET status_order = $1
	 WHERE id = $2`,
		statusOrder, orderID)

	if err != nil {
		return err
	}
	return nil
}
