package sddb

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type ProductsPostgresStorage struct {
	db *sqlx.DB
}

var errorsArticleNotFound = errors.New("no rows updated, article not found")

func NewProductsPostgresStorage(db *sqlx.DB) *ProductsPostgresStorage {
	return &ProductsPostgresStorage{db: db}
}

func (s *ProductsPostgresStorage) AddProduct(ctx context.Context, product Products) error {

	if _, err := s.db.ExecContext(
		ctx,
		`INSERT INTO products (article, catalog, name, description, photo_url, price, length, width, heigth, weight,available)
	    				VALUES ($1, $2, $3,$4, $5, $6,$7, $8, $9, $10,$11)
	    				ON CONFLICT DO NOTHING;`,
		product.Article,
		product.Catalog,
		product.Name,
		product.Description,
		pq.Array(product.PhotoUrl),
		product.Price,
		product.Length,
		product.Width,
		product.Height,
		product.Weight,
		product.Availability,
		//users.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (s *ProductsPostgresStorage) ChangeProductByArticle(ctx context.Context, product Products) error {

	result, err := s.db.ExecContext(
		ctx,
		`UPDATE products SET catalog = $2, name = $3, description = $4, photo_url = $5, price = $6, length = $7, width = $8, heigth = $9, weight = $10, available = $11
	    				WHERE article = $1`,
		product.Article,
		product.Catalog,
		product.Name,
		product.Description,
		pq.Array(product.PhotoUrl),
		product.Price,
		product.Length,
		product.Width,
		product.Height,
		product.Weight,
		product.Availability,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errorsArticleNotFound
	}

	return nil
}

func (s *ProductsPostgresStorage) GetCatalogNamesIsAvailable(ctx context.Context) ([]string, error) {
	var catalog []string

	if err := s.db.SelectContext(ctx, &catalog, `SELECT DISTINCT catalog FROM products where available = true`); err != nil {
		return nil, err
	}

	return catalog, nil
}

func (s *ProductsPostgresStorage) GetCatalogNames(ctx context.Context) ([]string, error) {

	var catalog []string
	if err := s.db.SelectContext(ctx, &catalog, `SELECT DISTINCT catalog FROM products`); err != nil {
		return nil, err
	}

	return catalog, nil
}

func (s *ProductsPostgresStorage) GetAllProducts(ctx context.Context) ([]Products, error) {

	var getProducts []getdbProduct
	if err := s.db.SelectContext(ctx,
		&getProducts,
		`SELECT
     article AS a_article,
     catalog AS a_catalog,
     name AS a_name,
     description AS a_description,
     photo_url AS a_photo_url,
     price AS a_price,
     length AS a_length,
     width  AS a_width,
     heigth  AS a_height,
     weight  AS a_weight,
    available as a_available
	 FROM products
	 ORDER BY article asc;`,
	); err != nil {
		return nil, err
	}

	var products []Products

	for _, getProduct := range getProducts {

		var photoUrls [][]byte
		for _, byteData := range getProduct.PhotoUrl {
			photoUrls = append(photoUrls, []byte(byteData))
		}

		product := Products{
			Article:      getProduct.Article,
			Catalog:      getProduct.Catalog,
			Name:         getProduct.Name,
			Description:  getProduct.Description,
			PhotoUrl:     photoUrls,
			Price:        getProduct.Price,
			Length:       getProduct.Length,
			Width:        getProduct.Width,
			Height:       getProduct.Height,
			Weight:       getProduct.Weight,
			Availability: getProduct.Availability,
		}

		products = append(products, product)
	}

	return products, nil
}

func (s *ProductsPostgresStorage) GetProductsByCatalogIsAvailable(ctx context.Context, ctlg string) ([]Products, error) {

	var getProducts []getdbProduct

	if err := s.db.SelectContext(ctx,
		&getProducts,
		`SELECT
       article AS a_article,
     catalog AS a_catalog,
     name AS a_name,
     description AS a_description,
     photo_url AS a_photo_url,
     price AS a_price,
     length AS a_length,
     width  AS a_width,
     heigth  AS a_height,
     weight  AS a_weight,
    available as a_available
	 FROM products
	 WHERE catalog = $1 and available = true`,
		ctlg); err != nil {
		return nil, err
	}

	var products []Products

	for _, getProduct := range getProducts {

		var photoUrls [][]byte
		for _, byteData := range getProduct.PhotoUrl {
			photoUrls = append(photoUrls, []byte(byteData))
		}

		product := Products{
			Article:      getProduct.Article,
			Catalog:      getProduct.Catalog,
			Name:         getProduct.Name,
			Description:  getProduct.Description,
			PhotoUrl:     photoUrls,
			Price:        getProduct.Price,
			Length:       getProduct.Length,
			Width:        getProduct.Width,
			Height:       getProduct.Height,
			Weight:       getProduct.Weight,
			Availability: getProduct.Availability,
		}

		products = append(products, product)
	}

	return products, nil
}

func (s *ProductsPostgresStorage) GetProductsByCatalog(ctx context.Context, ctlg string) ([]Products, error) {

	//var products []Products
	var getProducts []getdbProduct
	if err := s.db.SelectContext(ctx,
		&getProducts,
		`SELECT
     article AS a_article,
     catalog AS a_catalog,
     name AS a_name,
     description AS a_description,
     photo_url AS a_photo_url,
     price AS a_price
	 FROM products
	 WHERE catalog = $1`,
		ctlg); err != nil {
		return nil, err
	}

	var products []Products

	for _, getProduct := range getProducts {

		var photoUrls [][]byte
		for _, byteData := range getProduct.PhotoUrl {
			photoUrls = append(photoUrls, []byte(byteData))
		}

		product := Products{
			Article:     getProduct.Article,
			Catalog:     getProduct.Catalog,
			Name:        getProduct.Name,
			Description: getProduct.Description,
			PhotoUrl:    photoUrls,
			Price:       getProduct.Price,
			Length:      getProduct.Length,
			Width:       getProduct.Width,
			Height:      getProduct.Height,
			Weight:      getProduct.Weight,
		}

		products = append(products, product)
	}

	return products, nil
}

func (s *ProductsPostgresStorage) GetProductByArticle(ctx context.Context, article int) (Products, error) {

	var products dbProduct
	row := s.db.QueryRowContext(ctx,
		`SELECT
     			article AS c_article,
     			catalog AS c_catalog,
     			name AS c_name,
     			description AS c_description,
     			photo_url AS c_photo_url,
     			price AS c_price,
     			length AS c_length,
     			width AS c_width,
     			heigth AS c_height,
     			weight AS c_weight
	 			FROM products
	 			WHERE (article = $1)`,
		article)

	var photoUrl pq.ByteaArray

	err := row.Scan(
		&products.Article,
		&products.Catalog,
		&products.Name,
		&products.Description,
		&photoUrl,
		&products.Price,
		&products.Length,
		&products.Width,
		&products.Height,
		&products.Weight)

	if err != nil {
		return Products{}, err
	}

	products.PhotoUrl = photoUrl

	return Products{
		Article:     products.Article,
		Catalog:     products.Catalog,
		Name:        products.Name,
		Description: products.Description,
		PhotoUrl:    products.PhotoUrl,
		Price:       products.Price,
		Length:      products.Length,
		Width:       products.Width,
		Height:      products.Height,
		Weight:      products.Weight,
	}, err

}

type dbProduct struct {
	Article      int      `db:"a_article"`
	Catalog      string   `db:"a_catalog"`
	Name         string   `db:"a_name"`
	Description  string   `db:"a_description"`
	PhotoUrl     [][]byte `db:"a_photo_url"`
	Price        float64  `db:"a_price"`
	Length       int      `db:"a_length"`
	Width        int      `db:"a_width"`
	Height       int      `db:"a_height"`
	Weight       int      `db:"a_weight"`
	Availability bool     `db:"a_available"`
}
type getdbProduct struct {
	Article      int           `db:"a_article"`
	Catalog      string        `db:"a_catalog"`
	Name         string        `db:"a_name"`
	Description  string        `db:"a_description"`
	PhotoUrl     pq.ByteaArray `db:"a_photo_url"`
	Price        float64       `db:"a_price"`
	Length       int           `db:"a_length"`
	Width        int           `db:"a_width"`
	Height       int           `db:"a_height"`
	Weight       int           `db:"a_weight"`
	Availability bool          `db:"a_available"`
}
