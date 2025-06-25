package workflows

import (
	"fmt"
	"time"

	ordersv1 "github.com/yaninyzwitty/temporal-microservice-go/gen/orders/v1"
	"github.com/yaninyzwitty/temporal-microservice-go/services/order-service/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// CreateOrderWorkflow is the temporal workflow that CheckCustomerExists, CheckProductsAvailability and ReserveStock
func CreateOrderWorkflow(ctx workflow.Context, order *ordersv1.Order) error {

	// Define the activity options, including the retry policy

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second, //amount of time that must elapse before the first retry occurs
			MaximumInterval:    time.Minute, //maximum interval between retries
			BackoffCoefficient: 2,           //how much the retry interval increases
			// MaximumAttempts: 5, // Uncomment this if you want to limit attempts,
			NonRetryableErrorTypes: []string{},
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)
	var orderActivityClient *activities.OrderActivity

	// first check if customer exists

	var exists bool
	err := workflow.ExecuteActivity(ctx, orderActivityClient.CheckCustomerExists, order.CustomerId).Get(ctx, &exists)

	if err != nil {
		return fmt.Errorf("failed to check customer exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("customer not found")
	}

	// check if products are available
	var productExists bool
	err = workflow.ExecuteActivity(ctx, orderActivityClient.CheckProductsAvailability, order.Items).Get(ctx, &productExists)
	if err != nil {
		return fmt.Errorf("failed to check products availability: %w", err)
	}

	if !productExists {
		return fmt.Errorf("products not found")
	}

	// reserve stock
	err = workflow.ExecuteActivity(ctx, orderActivityClient.ReserveStock, order.Items).Get(ctx, &productExists)

	if err != nil {
		return fmt.Errorf("failed to reserve stock: %w", err)
	}

	// create order
	err = workflow.ExecuteActivity(ctx, orderActivityClient.CreateOrder, order.Items, order.OrderId, order.CustomerId, order.Status).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil

}
