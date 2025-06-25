package repository

import (
	"context"
	"fmt"

	ordersv1 "github.com/yaninyzwitty/temporal-microservice-go/gen/orders/v1"
	"github.com/yaninyzwitty/temporal-microservice-go/services/order-service/workflows"
	"go.temporal.io/sdk/client"
)

type OrderRepository struct {
	client client.Client
}

func NewOrderRepository(client client.Client) *OrderRepository {
	return &OrderRepository{
		client: client,
	}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *ordersv1.Order) error {

	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("order-%d", order.OrderId),
		TaskQueue: "order-service",
	}

	we, err := r.client.ExecuteWorkflow(ctx, workflowOptions, workflows.CreateOrderWorkflow, order)
	if err != nil {
		return fmt.Errorf("failed to execute workflow: %w", err)
	}

	// Wait for workflow completion and get any error result
	err = we.Get(ctx, nil)
	if err != nil {
		return fmt.Errorf("workflow execution failed: %w", err)
	}

	return nil
}

func (r *OrderRepository) GetOrder(orderId int64) (*ordersv1.Order, error) {
	return nil, nil
}

func (r *OrderRepository) UpdateOrder(order *ordersv1.Order) (*ordersv1.Order, error) {
	return nil, nil
}
