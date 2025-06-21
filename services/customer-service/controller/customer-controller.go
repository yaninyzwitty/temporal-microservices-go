package controller

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	v1 "github.com/yaninyzwitty/temporal-microservice-go/gen/customers/v1"
	"github.com/yaninyzwitty/temporal-microservice-go/gen/customers/v1/customersv1connect"
	"github.com/yaninyzwitty/temporal-microservice-go/services/customer-service/repository"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/snowflake"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CustomerController struct {
	customersv1connect.UnimplementedCustomersServiceHandler
	customerRepository *repository.CustomerRepository
}

func NewCustomerController(customerRepository *repository.CustomerRepository) *CustomerController {
	return &CustomerController{
		customerRepository: customerRepository,
	}
}

func (c *CustomerController) CreateCustomer(ctx context.Context, req *connect.Request[v1.CreateCustomerRequest]) (*connect.Response[v1.CreateCustomerResponse], error) {
	if req.Msg.Username == "" || req.Msg.AliasName == "" || req.Msg.Email == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("username, alias name and email are required"))
	}

	customerId, err := snowflake.GenerateID()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	customer := &v1.Customer{
		Id:        int64(customerId),
		Username:  req.Msg.Username,
		AliasName: req.Msg.AliasName,
		Email:     req.Msg.Email,
		CreatedAt: timestamppb.New(time.Now()),
	}

	if err := c.customerRepository.CreateCustomer(customer); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[v1.CreateCustomerResponse]{
		Msg: &v1.CreateCustomerResponse{
			Customer: customer,
		},
	}, nil
}
func (c *CustomerController) GetCustomer(ctx context.Context, req *connect.Request[v1.GetCustomerRequest]) (*connect.Response[v1.GetCustomerResponse], error) {
	if req.Msg.Id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}

	customer, err := c.customerRepository.GetCustomer(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[v1.GetCustomerResponse]{
		Msg: &v1.GetCustomerResponse{
			Customer: customer,
		},
	}, nil
}
func (c *CustomerController) DeleteCustomer(ctx context.Context, req *connect.Request[v1.DeleteCustomerRequest]) (*connect.Response[v1.DeleteCustomerResponse], error) {
	if req.Msg.Id <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}

	if err := c.customerRepository.DeleteCustomer(req.Msg.Id); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return &connect.Response[v1.DeleteCustomerResponse]{
		Msg: &v1.DeleteCustomerResponse{
			Success: true,
		},
	}, nil
}
