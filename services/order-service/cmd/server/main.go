package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/yaninyzwitty/temporal-microservice-go/gen/orders/v1/ordersv1connect"
	"github.com/yaninyzwitty/temporal-microservice-go/services/order-service/cmd/controller"
	"github.com/yaninyzwitty/temporal-microservice-go/services/order-service/repository"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg"
	database "github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/db"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/helpers"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/snowflake"
	"go.temporal.io/sdk/client"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	// implement order service

	cfg := pkg.Config{}

	if err := cfg.LoadConfig("config.yaml"); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}
	if err := snowflake.InitSonyFlake(); err != nil {
		slog.Error("failed to initialize snowflake", "error", err)
		os.Exit(1)
	}

	if err := godotenv.Load(); err != nil {
		slog.Error("failed to load .env file", "error", err)
		os.Exit(1)
	}

	astraCfg := &database.AstraConfig{
		Username: cfg.Database.Username,
		Path:     cfg.Database.Path,
		Token:    helpers.GetEnvOrDefault("ASTRA_TOKEN", ""),
	}

	db := database.NewAstraDB()
	session, err := db.Connect(context.Background(), astraCfg, 30*time.Second)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer session.Close()

	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		slog.Error("Unable to create client", "error", err)
		os.Exit(1)
	}

	defer temporalClient.Close()

	orderServiceAddr := fmt.Sprintf("localhost:%d", cfg.ProductServer.Port)
	orderRepository := repository.NewOrderRepository(temporalClient)
	orderController := controller.NewOrderController(orderRepository)

	mux := http.NewServeMux()

	orderPath, orderHandler := ordersv1connect.NewOrderServiceHandler(orderController)
	mux.Handle(orderPath, orderHandler)

	server := &http.Server{
		Addr:    orderServiceAddr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-quit
		slog.Info("received shutdown signal", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("server forced to shutdown", "error", err)
		} else {
			slog.Info("server shutdown gracefully")
		}
	}()

	// start http server
	slog.Info("starting HTTP server", "address", orderServiceAddr, "pid", os.Getpid())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}

	wg.Wait()
	slog.Info("service shutdown complete")
}
