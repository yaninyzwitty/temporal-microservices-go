package activities

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	ordersv1 "github.com/yaninyzwitty/temporal-microservice-go/gen/orders/v1"
)

type OrderActivity struct {
	Cassandra *gocql.Session
}

// ✅ Check if customer exists
func (o *OrderActivity) CheckCustomerExists(ctx context.Context, customerID int64) (bool, error) {
	var name string
	query := `SELECT name FROM customers WHERE id = ?`

	err := o.Cassandra.Query(query, customerID).WithContext(ctx).Scan(&name)
	if err != nil {
		return false, err
	}

	return name != "", nil
}

// ✅ Check products availability
func (o *OrderActivity) CheckProductsAvailability(ctx context.Context, items []*ordersv1.OrderItem) error {
	for _, item := range items {
		var stock int
		query := `SELECT stock FROM products WHERE id = ? LIMIT 1`

		err := o.Cassandra.Query(query, item.ProductId).WithContext(ctx).Scan(&stock)
		if err != nil {
			return fmt.Errorf("product %d not found", item.ProductId)
		}

		if stock < int(item.Quantity) {
			return fmt.Errorf("product %d has insufficient stock", item.ProductId)
		}
	}
	return nil
}

// ✅ Reserve stock atomically
func (o *OrderActivity) ReserveStock(ctx context.Context, items []*ordersv1.OrderItem) error {
	for _, item := range items {
		query := `UPDATE products SET stock = stock - ? WHERE id = ? IF stock >= ?`
		applied, err := o.Cassandra.Query(query, item.Quantity, item.ProductId, item.Quantity).WithContext(ctx).ScanCAS()

		if err != nil {
			return fmt.Errorf("failed reserving stock for product %d: %w", item.ProductId, err)
		}

		if !applied {
			return fmt.Errorf("insufficient stock for product %d", item.ProductId)
		}
	}
	return nil
}

func (o *OrderActivity) CreateOrder(ctx context.Context, items []*ordersv1.OrderItem, orderId, customerId int64, status string) error {
	createdAt := time.Now()
	updatedAt := createdAt

	// ✅ Insert into orders table

	orderQuery := `INSERT INTO orders (id, customer_id, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`

	err := o.Cassandra.Query(orderQuery, orderId, customerId, status, createdAt, updatedAt).WithContext(ctx).Exec()
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	// ✅ Insert into order_items table

	for _, item := range items {
		itemQuery := `INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)`

		if err := o.Cassandra.Query(itemQuery,
			orderId,
			item.ProductId,
			item.Quantity,
			item.Price,
		).WithContext(ctx).Exec(); err != nil {
			return fmt.Errorf("failed to insert order item (product_id=%d): %w", item.ProductId, err)
		}
	}

	return nil
}
