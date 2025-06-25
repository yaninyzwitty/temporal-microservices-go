// service worker to execute workflows and activities
package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/yaninyzwitty/temporal-microservice-go/services/order-service/activities"
	"github.com/yaninyzwitty/temporal-microservice-go/services/order-service/workflows"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg"
	database "github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/db"
	"github.com/yaninyzwitty/temporal-microservice-go/shared/pkg/helpers"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	cfg := pkg.Config{}

	if err := cfg.LoadConfig("config.yaml"); err != nil {
		slog.Error("failed to load config", "error", err)
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
	// Create the Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		slog.Error("Unable to create Temporal client", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	// Create the Temporal worker
	w := worker.New(c, "order-service-queue", worker.Options{})

	// inject cassandra session to orderactivity struct
	orderActivities := activities.OrderActivity{
		Cassandra: session,
	}

	// Register the workflow functions
	w.RegisterWorkflow(workflows.CreateOrderWorkflow)

	// register activities
	w.RegisterActivity(orderActivities)
	// run the worker
	if err := w.Run(worker.InterruptCh()); err != nil {
		slog.Error("Unable to start temporal worker", "error", err)
		os.Exit(1)
	}

}
