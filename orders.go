package sddb

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
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

func (s *OrdersPostgresStorage) GetOrdersByTgID(ctx context.Context, tgId int64) ([]Orders, error) {

	var corzine []Orders

	if err := s.db.SelectContext(ctx,
		&corzine,
		`SELECT
     			id ,
     			tg_id ,
	 			user_name ,
	 			first_name ,
	 			last_name ,
     			status_order ,
     			pvz ,
     			type_dostavka ,
     			orderr,
     			created_at  ,
     			update_at 
	 			FROM orders
	 			WHERE tg_id = $1`,
		tgId); err != nil {
		return nil, err
	}

	return corzine, nil
}

// use tg
func (s *OrdersPostgresStorage) AddOrder(ctx context.Context, order Orders) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO orders (tg_id,user_name,first_name,last_name,status_order,  orderr, update_at, pvz, type_dostavka)
	    				VALUES ($1, $2, $3, $4,$5,$6,$7,$8,$9)
	    				ON CONFLICT DO NOTHING;`,
		order.TgID,
		order.UserName,
		order.FirstName,
		order.LastName,
		order.StatusOrder,
		//
		order.Order,
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

	var corzine []Orders
	if err := s.db.SelectContext(ctx,
		&corzine,
		`SELECT
     			id,
     			tg_id,
     			user_name,
     			first_name,
	 			last_name,
     			status_order,
     			pvz,
     			type_dostavka,
     			orderr,
     			created_at ,
     			update_at
	 			FROM orders
	 			WHERE status_order = $1`,
		statusOrder); err != nil {
		return nil, err
	}
	return corzine, nil
}

func (s *OrdersPostgresStorage) GetOrderByTimeAndStatus(ctx context.Context, statusOrder int, time2 time.Time) ([]Orders, error) {

	var corzine []Orders

	if err := s.db.SelectContext(ctx,
		&corzine,
		`SELECT
     			id ,
     			tg_id ,
     			user_name,
     			first_name ,
	 			last_name ,
     			status_order ,
     			pvz ,
     			type_dostavka ,
     			orderr ,
     			created_at  AT TIME ZONE 'UTC' ,
     				update_at  AT TIME ZONE 'UTC' 
	 			FROM orders
	 		    WHERE status_order = $1 AND created_at > $2`,
		statusOrder, time2); err != nil {
		return nil, err
	}

	return corzine, nil
}
func (s *OrdersPostgresStorage) UpdateOrderByStatus(ctx context.Context, statusOrder int, orderID int) error {

	_, err := s.db.ExecContext(ctx,
		`UPDATE orders
	 SET status_order = $1
	 WHERE id = $2`,
		statusOrder, orderID)

	if err != nil {
		return err
	}
	return nil
}
