package pkg

import (
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CustomerServer CustomerServer `yaml:"customer_server"`
	Database       Database       `yaml:"database"`
	ProductServer  ProductServer  `yaml:"products-server"`
}

type Database struct {
	Username string `yaml:"username"`
	// Token           string        `yaml:"token"` -- use .env
	Path            string   `yaml:"path"`
	Keyspace        string   `yaml:"keyspace"`
	Hosts           []string `yaml:"hosts"`
	LocalDataCenter string   `yaml:"localDataCenter"`
	// MaxRetries      int           `yaml:"maxRetries"`
	// Timeout time.Duration `yaml:"timeout"`
}

type CustomerServer struct {
	Port int `yaml:"port"`
}

type ProductServer struct {
	Port int `yaml:"port"`
}

func (c *Config) LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("failed to read config file", "error", err, "path", path)
		return err
	}

	err = yaml.Unmarshal(data, c)
	if err != nil {
		slog.Error("failed to unmarshal config file", "error", err, "path", path)
		return err
	}

	return nil

}
