package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	"connectrpc.com/connect"
	v1 "github.com/yaninyzwitty/temporal-microservice-go/gen/orders/v1"
	"github.com/yaninyzwitty/temporal-microservice-go/gen/orders/v1/ordersv1connect"
	"github.com/yaninyzwitty/temporal-microservice-go/services/order-service/repository"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/snowflake"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type OrderController struct {
	ordersv1connect.UnimplementedOrderServiceHandler
	orderRepository *repository.OrderRepository
}

func NewOrderController(orderRepository *repository.OrderRepository) *OrderController {
	return &OrderController{
		orderRepository: orderRepository,
	}
}

func (c *OrderController) CreateOrder(ctx context.Context, req *connect.Request[v1.CreateOrderRequest]) (*connect.Response[v1.CreateOrderResponse], error) {
	if req.Msg.CustomerId <= 0 || len(req.Msg.Items) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("customer_id and items are required"))
	}

	orderId, err := snowflake.GenerateID()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	now := time.Now()

	// create order obj
	order := &v1.Order{
		OrderId:    int64(orderId),
		CustomerId: req.Msg.CustomerId,
		Items:      req.Msg.Items,
		Status:     v1.OrderStatus_ORDER_STATUS_CREATED,
		CreatedAt:  timestamppb.New(now),
		UpdatedAt:  nil,
	}

	// save order to db
	if err := c.orderRepository.CreateOrder(ctx, order); err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed saving order: %w", err))
	}

	// build the response
	resp := &v1.CreateOrderResponse{
		Order: order,
	}

	return connect.NewResponse(resp), nil

}
