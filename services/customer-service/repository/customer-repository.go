package repository

import (
	"time"

	"github.com/gocql/gocql"
	customersv1 "github.com/yaninyzwitty/temporal-microservice-go/gen/customers/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CustomerRepository struct {
	session *gocql.Session
}

func NewCustomerRepository(session *gocql.Session) *CustomerRepository {
	return &CustomerRepository{
		session: session,
	}
}

func (r *CustomerRepository) CreateCustomer(customer *customersv1.Customer) error {
	query := `
		INSERT INTO customers (id, username, alias_name, email, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	return r.session.Query(query, customer.Id, customer.Username, customer.AliasName, customer.Email, customer.CreatedAt).Exec()

}

func (r *CustomerRepository) GetCustomer(id int64) (*customersv1.Customer, error) {
	var createdAt time.Time
	query := `
		SELECT id, username, alias_name, email, created_at
		FROM customers
		WHERE id = ?
	`

	var customer customersv1.Customer
	if err := r.session.Query(query, id).Scan(&customer.Id, &customer.Username, &customer.AliasName, &customer.Email, &createdAt); err != nil {
		return nil, err
	}

	customer.CreatedAt = timestamppb.New(createdAt)

	return &customer, nil

}

func (r *CustomerRepository) DeleteCustomer(id int64) error {
	query := `
		DELETE FROM customers
		WHERE id = ?
	`

	return r.session.Query(query, id).Exec()

}
