# Temporal Microservice Go

A distributed microservices architecture built with Go, Temporal, and ScyllaDB, implementing an e-commerce system with customer, product, and order management.

## 🏗️ Architecture

This project implements a microservices architecture with the following components:

- **Customer Service**: Manages customer data and operations
- **Product Service**: Handles product catalog and inventory
- **Order Service**: Processes orders with Temporal workflows
- **Worker Service**: Executes Temporal workflows and activities
- **ScyllaDB**: Distributed NoSQL database for data persistence

## 🚀 Features

- **gRPC/Connect-RPC**: Modern RPC framework for service communication
- **Temporal Workflows**: Reliable, fault-tolerant business logic orchestration
- **ScyllaDB Integration**: High-performance distributed database
- **Protocol Buffers**: Type-safe API definitions
- **Docker Compose**: Local development environment
- **Graceful Shutdown**: Proper service lifecycle management

## 📁 Project Structure

```
temporal-microservice-go/
├── proto/                    # Protocol Buffer definitions
│   ├── customers/v1/
│   ├── products/v1/
│   └── orders/v1/
├── gen/                      # Generated Go code from protobuf
├── services/                 # Microservices
│   ├── customer-service/
│   ├── product-service/
│   ├── order-service/
│   └── worker/
├── shared/                   # Shared packages
│   └── pkg/
│       ├── config.go
│       ├── db/
│       ├── helpers/
│       └── snowflake/
├── schema/                   # Database schemas
├── compose.yaml             # Docker Compose configuration
├── config.yaml              # Application configuration
├── buf.yaml                 # Buf configuration
└── buf.gen.yaml            # Buf code generation config
```

## 🛠️ Prerequisites

- Go 1.24.3+
- Docker and Docker Compose
- Buf CLI
- Temporal Server (for local development)

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <repository-url>
cd temporal-microservice-go
```

### 2. Start Infrastructure

```bash
# Start ScyllaDB cluster
docker-compose up -d

# Start Temporal Server (if not already running)
temporal server start-dev
```

### 3. Generate Protocol Buffer Code

```bash
# Install buf CLI if not already installed
go install github.com/bufbuild/buf/cmd/buf@latest

# Generate Go code from protobuf definitions
buf generate
```

### 4. Set Up Environment Variables

Create a `.env` file in the root directory:

```env
ASTRA_TOKEN=your_astra_token_here
```

### 5. Run the Services

```bash
# Terminal 1: Start Customer Service
go run services/customer-service/cmd/server/main.go

# Terminal 2: Start Product Service
go run services/product-service/cmd/server/main.go

# Terminal 3: Start Order Service
go run services/order-service/cmd/server/main.go

# Terminal 4: Start Worker Service
go run services/worker/cmd/server/main.go
```

## 📋 Configuration

The application uses `config.yaml` for configuration:

```yaml
customer_server:
  port: 50051
products-server:
  port: 50052
order-server:
  port: 50053
database:
  username: token
  token: token
  path: ./secure-connect.zip
  keyspace: temporal_microservice_keyspace
  hosts:
    - "127.0.0.1:9042"
    - "127.0.0.1:9043"
    - "127.0.0.1:9044"
  localDataCenter: "scylla-net"
```

## 🔧 Development

### Protocol Buffer Development

1. **Edit Protobuf Files**: Modify files in `proto/` directory
2. **Generate Code**: Run `buf generate` to update generated Go code
3. **Lint**: Run `buf lint` to check for issues
4. **Breaking Changes**: Run `buf breaking` to detect breaking changes

### Database Schema

The project includes CQL schemas for ScyllaDB:

- Customer service schema: `schema/customers.cql`
- Additional schemas can be added for products and orders

### Adding New Services

1. Create a new directory in `services/`
2. Define protobuf messages and services in `proto/`
3. Generate code with `buf generate`
4. Implement the service following the existing patterns

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific service tests
go test ./services/customer-service/...
```

## 📊 Monitoring and Observability

- **Structured Logging**: Uses `slog` for consistent logging
- **Graceful Shutdown**: Proper signal handling for clean service termination
- **Error Handling**: Comprehensive error handling with retry policies

## 🔒 Security

- **Environment Variables**: Sensitive data stored in `.env` files
- **Database Security**: Uses Astra DB with secure connections
- **Input Validation**: Protobuf-based type safety

## 🚀 Deployment

### Docker Deployment

```bash
# Build all services
docker build -t temporal-microservice-go .

# Run with docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes Deployment

Kubernetes manifests can be added to deploy the services in a Kubernetes cluster.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Buf Generation Fails**: Ensure buf CLI is installed and up to date
2. **Database Connection Issues**: Check ScyllaDB is running and accessible
3. **Temporal Connection Issues**: Verify Temporal server is running
4. **Port Conflicts**: Check if required ports are available

### Logs

Check service logs for detailed error information:

```bash
# Customer service logs
go run services/customer-service/cmd/server/main.go 2>&1 | tee customer.log

# Worker service logs
go run services/worker/cmd/server/main.go 2>&1 | tee worker.log
```

## 📚 Additional Resources

- [Temporal Documentation](https://docs.temporal.io/)
- [ScyllaDB Documentation](https://docs.scylladb.com/)
- [Connect-RPC Documentation](https://connectrpc.com/)
- [Protocol Buffers Guide](https://developers.google.com/protocol-buffers)
