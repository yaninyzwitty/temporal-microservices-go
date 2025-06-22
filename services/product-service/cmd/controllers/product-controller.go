package controllers

import (
	"context"
	"errors"
	"strconv"
	"time"

	"connectrpc.com/connect"
	v1 "github.com/yaninyzwitty/temporal-microservice-go/gen/products/v1"
	v1connect "github.com/yaninyzwitty/temporal-microservice-go/gen/products/v1/v1connect"
	"github.com/yaninyzwitty/temporal-microservice-go/services/product-service/cmd/repository"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/snowflake"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProductController struct {
	v1connect.UnimplementedProductServiceHandler
	productRepository *repository.ProductRepository
}

func NewProductController(productRepository *repository.ProductRepository) *ProductController {
	return &ProductController{
		productRepository: productRepository,
	}
}

func (c *ProductController) CreateProduct(ctx context.Context, req *connect.Request[v1.CreateProductRequest]) (*connect.Response[v1.CreateProductResponse], error) {
	if req.Msg.Name == "" || req.Msg.Description == "" || req.Msg.Price == 0 || req.Msg.Currency == "" || req.Msg.ImageUrl == "" || req.Msg.Stock == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name, description, price, currency, image_url and stock are required"))
	}

	productId, err := snowflake.GenerateID()
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	product := &v1.Product{
		Id:          int64(productId),
		Name:        req.Msg.Name,
		Description: req.Msg.Description,
		Price:       req.Msg.Price,
		Currency:    req.Msg.Currency,
		ImageUrl:    req.Msg.ImageUrl,
		Stock:       req.Msg.Stock,
		CreatedAt:   timestamppb.New(time.Now()),
		UpdatedAt:   timestamppb.New(time.Now()),
	}

	if err := c.productRepository.CreateProduct(product); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[v1.CreateProductResponse]{
		Msg: &v1.CreateProductResponse{
			Product: product,
		},
	}, nil
}

func (c *ProductController) GetProduct(ctx context.Context, req *connect.Request[v1.GetProductRequest]) (*connect.Response[v1.GetProductResponse], error) {
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}

	productId, err := strconv.ParseUint(req.Msg.Id, 10, 64)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("cannot parse id to int64"))
	}

	product, err := c.productRepository.GetProduct(int64(productId))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[v1.GetProductResponse]{
		Msg: &v1.GetProductResponse{
			Product: product,
		},
	}, nil
}

func (c *ProductController) DeleteProduct(ctx context.Context, req *connect.Request[v1.DeleteProductRequest]) (*connect.Response[v1.DeleteProductResponse], error) {
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("id is required"))
	}

	productId, err := strconv.ParseUint(req.Msg.Id, 10, 64)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("cannot parse id to int64"))
	}

	if err := c.productRepository.DeleteProduct(int64(productId)); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return &connect.Response[v1.DeleteProductResponse]{
		Msg: &v1.DeleteProductResponse{
			Deleted: true,
		},
	}, nil
}
