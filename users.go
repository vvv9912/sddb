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

	if _, err := s.db.ExecContext(
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
	var status int
	var state int

	row := s.db.QueryRowContext(ctx, `SELECT status_user,state_user FROM users where tg_id = $1`, tgID)

	err := row.Scan(&status, &state)
	if err != nil {
		return 0, 0, err
	}

	return status, state, nil
}
func (s *UsersPostgresStorage) UpdateStateByTgID(ctx context.Context, tgId int64, state int) error {
	if _, err := s.db.ExecContext(
		ctx,
		`UPDATE users SET state_user = $1 WHERE (tg_id = $2)`,
		state,
		tgId,
	); err != nil {
		return err
	}

	return nil
}
