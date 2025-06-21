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

	"github.com/yaninyzwitty/temporal-microservice-go/gen/customers/v1/customersv1connect"
	"github.com/yaninyzwitty/temporal-microservice-go/services/customer-service/controller"
	"github.com/yaninyzwitty/temporal-microservice-go/services/customer-service/repository"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/db"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/helpers"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/snowflake"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	config := pkg.Config{}
	if err := config.LoadConfig("config.yaml"); err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	if err := snowflake.InitSonyFlake(); err != nil {
		slog.Error("failed to initialize snowflake", "error", err)
		os.Exit(1)
	}

	session, err := db.NewCassandra(&db.CassandraConfig{
		Hosts:      config.Database.Hosts,
		Keyspace:   config.Database.Keyspace,
		MaxRetries: 30,
		Username:   config.Database.Username,
		Token:      helpers.GetEnvOrDefault("ASTRA_DB_TOKEN", ""),
		Path:       config.Database.Path,
		Timeout:    30 * time.Second,
	})

	if err != nil {
		slog.Error("failed to create cassandra client", "error", err)
		os.Exit(1)
	}

	defer session.Close()

	customerRepository := repository.NewCustomerRepository(session)

	customerServiceAddr := fmt.Sprintf("localhost:%d", config.CustomerServer.Port)

	// Initialize controller and mux
	customerController := controller.NewCustomerController(customerRepository)
	customersPath, customersHandler := customersv1connect.NewCustomersServiceHandler(customerController)

	mux := http.NewServeMux()
	mux.Handle(customersPath, customersHandler)

	// Setup HTTP server with h2c (HTTP/2 cleartext)
	server := &http.Server{
		Addr:    customerServiceAddr,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// Graceful shutdown
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

	// Start HTTP server
	slog.Info("starting HTTP server", "address", customerServiceAddr, "pid", os.Getpid())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}

	wg.Wait()
	slog.Info("service shutdown complete")
}
