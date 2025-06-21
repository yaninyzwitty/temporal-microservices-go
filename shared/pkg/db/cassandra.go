package db

import (
	"fmt"
	"log/slog"
	"time"

	gocqlastra "github.com/datastax/gocql-astra"
	"github.com/gocql/gocql"
)

type CassandraConfig struct {
	Hosts      []string
	Keyspace   string
	MaxRetries int
	Username   string
	Token      string
	Path       string
	Timeout    time.Duration
}

func NewCassandra(config *CassandraConfig) (*gocql.Session, error) {
	// If any of the cloud (Astra) fields are missing, assume local setup
	if config.Username == "" || config.Token == "" || config.Path == "" {
		var session *gocql.Session
		var err error

		for attempt := 1; attempt <= config.MaxRetries; attempt++ {
			cluster := gocql.NewCluster(config.Hosts...)
			cluster.Keyspace = config.Keyspace
			cluster.Consistency = gocql.Quorum
			cluster.Timeout = config.Timeout

			session, err = cluster.CreateSession()
			if err == nil {
				// Test connection
				if pingErr := session.Query("SELECT now() FROM system.local").Exec(); pingErr == nil {
					slog.Info("Connected to local Cassandra/Scylla")
					return session, nil
				}
				session.Close()
				err = fmt.Errorf("ping failed after session creation")
			}

			slog.Error("Failed to create local Cassandra session", "attempt", attempt, "error", err)
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		return nil, fmt.Errorf("exceeded max retries for local Cassandra connection")
	}

	// Use Astra DB setup
	cluster, err := gocqlastra.NewClusterFromBundle(config.Path, config.Username, config.Token, config.Timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Astra cluster: %w", err)
	}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create Astra Cassandra session: %w", err)
	}

	slog.Info("Connected to Astra Cassandra")
	return session, nil
}
