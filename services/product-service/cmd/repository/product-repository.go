package repository

import (
	"log/slog"
	"time"

	"github.com/gocql/gocql"
	v1 "github.com/yaninyzwitty/temporal-microservice-go/gen/products/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductRepository struct {
	session *gocql.Session
}

func NewProductRepository(session *gocql.Session) *ProductRepository {
	return &ProductRepository{session: session}
}

func (r *ProductRepository) CreateProduct(product *v1.Product) error {

	query := `
		INSERT INTO products_keyspace.products (id, name, description, price, currency, image_url, stock, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	return r.session.Query(query, product.Id, product.Name, product.Description, product.Price, product.Currency, product.ImageUrl, product.Stock, product.CreatedAt.AsTime(), product.UpdatedAt.AsTime()).Exec()

}

func (r *ProductRepository) GetProduct(id int64) (*v1.Product, error) {
	var product v1.Product
	query := `
		SELECT id, name, description, price, currency, image_url, stock, created_at, updated_at	
		FROM products_keyspace.products
		WHERE id = ?
	`
	var createdAt, updatedAt time.Time
	if err := r.session.Query(query, id).Scan(&product.Id, &product.Name, &product.Description, &product.Price, &product.Currency, &product.ImageUrl, &product.Stock, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	product.CreatedAt = timestamppb.New(createdAt)
	product.UpdatedAt = timestamppb.New(updatedAt)
	return &product, nil
}

func (r *ProductRepository) DeleteProduct(id int64) error {
	slog.Info("Deleting product", "id", id)
	query := `
		DELETE FROM products_keyspace.products WHERE id = ?
	`
	return r.session.Query(query, id).Exec()
}
