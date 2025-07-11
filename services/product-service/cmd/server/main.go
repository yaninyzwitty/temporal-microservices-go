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
	"github.com/yaninyzwitty/temporal-microservice-go/gen/products/v1/v1connect"
	"github.com/yaninyzwitty/temporal-microservice-go/services/product-service/controllers"
	"github.com/yaninyzwitty/temporal-microservice-go/services/product-service/repository"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg"
	database "github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/db"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/helpers"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/snowflake"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
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

	productServiceAddr := fmt.Sprintf("localhost:%d", cfg.ProductServer.Port)
	productRepository := repository.NewProductRepository(session)
	productController := controllers.NewProductController(productRepository)

	productPath, productHandler := v1connect.NewProductServiceHandler(productController)

	mux := http.NewServeMux()
	mux.Handle(productPath, productHandler)

	server := &http.Server{
		Addr:    productServiceAddr,
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

	// Start HTTP server
	slog.Info("starting HTTP server", "address", productServiceAddr, "pid", os.Getpid())
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}

	wg.Wait()
	slog.Info("service shutdown complete")
}
