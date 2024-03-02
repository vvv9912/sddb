package sddb

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type UsersPostgresStorage struct {
	db *sqlx.DB
}

func NewUsersPostgresStorage(db *sqlx.DB) *UsersPostgresStorage {
	return &UsersPostgresStorage{db: db}
}

func (s *UsersPostgresStorage) AddUser(ctx context.Context, users Users) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO users (tg_id, status_user,state_user)
	    				VALUES ($1, $2, $3)
	    				ON CONFLICT DO NOTHING;`,
		users.TgID,
		users.StatusUser,
		users.StateUser,
		//users.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *UsersPostgresStorage) GetStatusUserByTgID(ctx context.Context, tgID int64) (int, int, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return 0, 0, err
	}
	defer conn.Close()

	var status int
	var state int
	row := conn.QueryRowContext(ctx, `SELECT status_user,state_user FROM users where tg_id = $1`, tgID)
	err = row.Scan(&status, &state)
	if err != nil {
		return 0, 0, err
	}
	//if err := conn.SelectContext(ctx, &status, `SELECT status_user FROM users WHERE tg_id = $1`, tgID); err != nil {
	//	return 0, err
	//}

	return status, state, nil
}
func (s *UsersPostgresStorage) UpdateStateByTgID(ctx context.Context, tgId int64, state int) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := conn.ExecContext(
		ctx,
		`UPDATE users SET state_user = $1 WHERE (tg_id = $2)`,
		state,
		tgId,
	); err != nil {
		return err
	}
	return nil
}

//func (s *UsersPostgresStorage) UpdateShopCartByTgId(ctx context.Context, tgId int64, corzina []int64) error {
//	conn, err := s.db.Connx(ctx)
//	if err != nil {
//		return err
//	}
//	defer conn.Close()
//
//	if _, err := conn.ExecContext(
//		ctx,
//		`UPDATE users SET corzina = $1 WHERE tg_id = $2`,
//		pq.Array(corzina),
//		tgId,
//	); err != nil {
//		return err
//	}
//
//	return nil
//}
////////////////////////////////////
//func (s *UsersPostgresStorage) GetCorzinaByID(ctx context.Context, id int) ([]int64, error) {
//	//conn, err := s.db.Connx(ctx)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//defer conn.Close()
//	////var corz []int64
//	//var corz pq.Int64Array
//	//
//	//if err := conn.SelectContext(ctx, &corz, `SELECT corzina FROM users where id = $1`, id); err != nil {
//	//	return nil, err
//	//}
//	//
//	//return corz, nil
//	//////
//	conn, err := s.db.Connx(ctx)
//	if err != nil {
//		return nil, err
//	}
//	defer conn.Close()
//	var corz []int64
//
//	row := conn.QueryRowContext(ctx, `SELECT corzina FROM users where id = $1`, id)
//
//	err = row.Scan(pq.Array(&corz))
//	if err != nil {
//		return nil, err
//	}
//	return corz, nil
//}
///////////////////////////////
//func (s *UsersPostgresStorage) GetCorzinaByTgID(ctx context.Context, tgID int64) ([]int64, error) {
//	conn, err := s.db.Connx(ctx)
//	if err != nil {
//		return nil, err
//	}
//	defer conn.Close()
//	var corz []int64
//	row := conn.QueryRowContext(ctx, `SELECT corzina FROM users where tg_id = $1`, tgID)
//	err = row.Scan(pq.Array(&corz))
//	if err != nil {
//		return nil, err
//	}
//	return corz, nil
//}
////////////////////////////
//type dbOrders struct {
//	ID          int    `db:"id"`
//	IDUser      int    `db:"id_user"`
//	StatusOrder int    `db:"status_order"`
//	Pvz         string `db:"pvz"`
//	Order       string `db:"order"`
//}
