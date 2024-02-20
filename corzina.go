package sddb

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"time"
)

type CorzinaPostgresStorage struct {
	db *sqlx.DB
}

func NewCorzinaPostgresStorage(db *sqlx.DB) *CorzinaPostgresStorage {
	return &CorzinaPostgresStorage{db: db}
}
func (s *CorzinaPostgresStorage) AddCorzinas(ctx context.Context, corz Corzine) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO corzina (tg_id, article,quantity,CREATED_AT)
	    				VALUES ($1, $2, $3, $4)
	    				ON CONFLICT DO NOTHING;`,
		corz.TgId,
		corz.Article,
		corz.Quantity,
		corz.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}
func (s *CorzinaPostgresStorage) CorzinaByTgId(ctx context.Context, tgId int64) ([]Corzine, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	//var products []model.Products
	var corzine []dbCorzine
	if err := conn.SelectContext(ctx,
		&corzine,
		`SELECT
     			id AS c_id,
     			tg_id AS c_tg_id,
     			article AS c_article,
     			quantity AS c_quantity,
     			CREATED_AT AS c_CREATED_AT
	 			FROM corzina
	 			WHERE tg_id = $1`,
		tgId); err != nil {
		return nil, err
	}

	return lo.Map(corzine, func(corzin dbCorzine, _ int) Corzine { return Corzine(corzin) }), nil
}

func (s *CorzinaPostgresStorage) CorzinaByTgIdANDAtricle(ctx context.Context, tgId int64, article int) (Corzine, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return Corzine{}, err
	}
	defer conn.Close()
	var corzine dbCorzine
	row := conn.QueryRowContext(ctx,
		`SELECT
     			id AS c_id,
     			tg_id AS c_tg_id,
     			article AS c_article,
     			quantity AS c_quantity,
     			CREATED_AT AS c_CREATED_AT
	 			FROM corzina
	 			WHERE (tg_id = $1 and article = $2)`,
		tgId, article)
	err = row.Scan(&corzine.ID, &corzine.TgId, &corzine.Article,
		&corzine.Quantity, &corzine.CreatedAt)
	if err != nil {
		return Corzine{}, err
	}
	//return lo.Map(corzine, func(corzin dbCorzine, _ int) model.Corzine { return model.Corzine(corzin) }), nil
	return Corzine{ID: corzine.ID, TgId: corzine.TgId, Article: corzine.Article, Quantity: corzine.Quantity, CreatedAt: corzine.CreatedAt}, nil
}
func (s *CorzinaPostgresStorage) UpdateCorzinaByTgId(ctx context.Context, tgId int64, article int, quantity int) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := conn.ExecContext(
		ctx,
		`UPDATE corzina SET quantity = $1 WHERE (tg_id = $2 AND article = $3)`,
		quantity,
		tgId,
		article,
	); err != nil {
		return err
	}
	return nil
}

type dbCorzine struct {
	ID        int       `db:"c_id"`
	TgId      int64     `db:"c_tg_id"`
	Article   int       `db:"c_article"` //В наличии
	Quantity  int       `db:"c_quantity"`
	CreatedAt time.Time `db:"c_created_at"`
}

func (s *CorzinaPostgresStorage) CorzinaByTgIdwithCalalog(ctx context.Context, tgId int64) ([]DbCorzineCatalog, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	//var products []model.Products
	var corzineall []DbCorzineCatalog
	if err := conn.SelectContext(ctx,
		&corzineall,
		`SELECT corzina.id AS c_id, 
       		   corzina.tg_id AS c_tg_id,
       		   p.article AS c_article,
       		   corzina.quantity AS c_quantity,
       		   p.catalog AS c_catalog,
       		   p.name AS c_name,
			   p.price AS c_price,
			   p.length AS c_length,
			   p.width as c_width,
			   p.heigth as c_height,
			   p.weight as c_weight
			   FROM  corzina
			   LEFT JOIN public.products p on corzina.article = p.article  where tg_id = $1 GROUP BY id, p.article;`,
		tgId); err != nil {
		return nil, err
	}

	return corzineall, nil
}
func (s *CorzinaPostgresStorage) DeleteCorzinaByTgID(ctx context.Context, tgId int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	row, err := conn.QueryxContext(ctx, `DELETE FROM corzina WHERE tg_id = $1;`, tgId)
	defer row.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *CorzinaPostgresStorage) DeleteCorzinaByTgIDandArticle(ctx context.Context, tgId int64, article int) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	//var products []model.Products
	row, err := conn.QueryxContext(ctx, `DELETE FROM corzina WHERE (tg_id = $1 and article = $2);`, tgId, article)
	defer row.Close()
	if err != nil {
		return err
	}
	return nil
}

type DbCorzineCatalog struct {
	ID       int     `db:"c_id"`
	TgId     int64   `db:"c_tg_id"`
	Article  int     `db:"c_article"` //В наличии
	Quantity int     `db:"c_quantity"`
	Catalog  string  `db:"c_catalog,"`
	Name     string  `db:"c_name"`
	Price    float64 `db:"c_price"`
	Length   int     `db:"c_length"`
	Width    int     `db:"c_width"`
	Height   int     `db:"c_height"`
	Weight   int     `db:"c_weight"`
}
