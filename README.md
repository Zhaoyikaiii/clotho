<div align="center">
  <img src="assets/logo.png" alt="Clotho" width="512">

  <h1>Clotho: Distributed Task Orchestration & Workflow Engine</h1>

  <h3>Temporal-style reliable workflows • Idempotent execution • Pluggable memory layer</h3>

  <p>
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white" alt="Go">
    <img src="https://img.shields.io/badge/Architecture-Distributed-blue" alt="Architecture">
    <img src="https://img.shields.io/badge/license-MIT-green" alt="License">
    <br>
    <a href="https://clotho.ai"><img src="https://img.shields.io/badge/Website-clotho.ai-blue?style=flat&logo=google-chrome&logoColor=white" alt="Website"></a>
    <a href="https://github.com/Zhaoyikaiii/clotho"><img src="https://img.shields.io/badge/GitHub-Repository-black?style=flat&logo=github&logoColor=white" alt="GitHub"></a>
  </p>

</div>

---

## What is Clotho?

Clotho is a **distributed task orchestration and workflow engine** inspired by Temporal, built in Go. It provides reliable workflow execution with features like:

- **Workflow DAG**: Define complex business logic as directed acyclic graphs
- **Idempotent Execution**: Safe retry with exactly-once semantics
- **Retry & Backoff**: Configurable retry policies with exponential backoff
- **Timeouts**: Activity and workflow timeout management
- **Compensation**: Saga pattern support for distributed transactions
- **Observability**: Built-in tracing, metrics, and logging
- **Pluggable Memory Layer**: Vector databases (Milvus), RAG, event archiving

## Core Features

| Feature | Description |
|---------|-------------|
| **Workflow Engine** | Define workflows as DAGs with activities, conditional branches, and parallel execution |
| **Task Queue** | Durable task queues with at-least-once delivery |
| **Workers** | Scalable worker pools with automatic load balancing |
| **State Persistence** | Workflow state stored in durable storage (SQLite, PostgreSQL ready) |
| **Memory Providers** | Pluggable memory/retrieval layer (Milvus, Qdrant, Weaviate, custom RAG) |
| **Observability** | OpenTelemetry tracing, Prometheus metrics, structured logging |
| **Fault Tolerance** | Automatic retry, timeout handling, dead letter queues |

## Architecture Overview

```
+-----------------------------------------------------------------+
|                         Clotho Platform                         |
+-----------------------------------------------------------------+
|  +--------------+  +--------------+  +----------------------+|
|  |   API Server |  |  Web UI      |  |   gRPC/REST APIs     ||
|  |  (Gateway)   |  |  (Optional)  |  |                      ||
|  +--------------+  +--------------+  +----------------------+|
|                                                                 |
|  +-----------------------------------------------------------+|
|  |              Orchestrator / Control Plane                  ||
|  |  + Workflow scheduling    + Task routing                  ||
|  |  + State management     + Event handling                ||
|  +-----------------------------------------------------------+|
|                                                                 |
|  +--------------+  +--------------+  +--------------------+ |
|  |   Workers    |  |   Workers    |  |    Workers         | |
|  |  (Activity 1)|  |  (Activity 2)|  |   (Activity N)     | |
|  +--------------+  +--------------+  +--------------------+ |
|                                                                 |
|  +-----------------------------------------------------------+|
|  |                     Storage Layer                          ||
|  |  +------------+ +------------+ +-----------------------+   | |
|  |  |  Workflow  | |  Events    | |   Memory Providers   |   | |
|  |  |   State    | |   Archive   | | (Milvus/Qdrant/RAG) |   | |
|  |  +------------+ +------------+ +-----------------------+   | |
|  +-----------------------------------------------------------+|
+-----------------------------------------------------------------+
```

## Quickstart

### 1. Installation

```bash
# Download from releases
wget https://github.com/Zhaoyikaiii/clotho/releases/download/v0.1.0/clotho-linux-amd64
chmod +x clotho-linux-amd64
./clotho-linux-amd64 --help

# Or build from source
git clone https://github.com/Zhaoyikaiii/clotho.git
cd clotho
make build
./build/clotho --help
```

### 2. Configuration

Create `~/.clotho/config.yaml`:

```yaml
version: "1.0"

# Workflow engine configuration
engine:
  address: "localhost:6789"
  workers:
    count: 4
    poll_interval: "1s"

# Storage backend
storage:
  type: sqlite
  path: "~/.clotho/data/clotho.db"

# Memory/Retrieval layer
memory:
  providers:
    - type: milvus
      endpoint: "localhost:19530"
      collection: "workflow_memory"
    - type: rag
      embedding_model: "openai/text-embedding-3-small"

# Observability
observability:
  log_level: info
  metrics:
    enabled: true
    port: 9090
  tracing:
    enabled: true
    endpoint: "http://localhost:14268/api/traces"

# Model providers (for AI-enhanced workflows)
providers:
  - name: openai
    api_key: "${OPENAI_API_KEY}"
    api_base: "https://api.openai.com/v1"
```

### 3. Start the Server

```bash
# Start the orchestration server
clotho server start

# Or run in Docker
docker run -d \
  --name clotho \
  -p 6789:6789 \
  -p 9090:9090 \
  -v ~/.clotho:/root/.clotho \
  zhaoyikaiii/clotho:latest
```

### 4. Run a Sample Workflow

```bash
# Register and run a workflow
clotho workflow register --file ./examples/workflow-example.yaml
clotho workflow run --name order-processing --input '{"order_id": "12345"}'

# View workflow status
clotho workflow status --id <workflow-id>
```

## Configuration Reference

### Configuration File

Default location: `~/.clotho/config.yaml`

| Field | Type | Description | Default |
|-------|------|-------------|---------|
| `version` | string | Config version | "1.0" |
| `engine.address` | string | gRPC server address | "localhost:6789" |
| `engine.workers.count` | int | Number of workers | 4 |
| `engine.workers.poll_interval` | string | Task poll interval | "1s" |
| `storage.type` | string | Storage backend | "sqlite" |
| `storage.path` | string | Database path | "~/.clotho/data/clotho.db" |
| `memory.providers` | array | Memory providers | [] |
| `observability.log_level` | string | Log level | "info" |
| `observability.metrics.enabled` | bool | Enable metrics | true |
| `observability.tracing.enabled` | bool | Enable tracing | false |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `CLOTHO_HOME` | Config directory (~/.clotho) |
| `CLOTHO_CONFIG` | Config file path |
| `CLOTHO_LOG_LEVEL` | Log level (debug/info/warn/error) |
| `CLOTHO_WORKERS` | Number of workers |
| `CLOTHO_STORAGE_TYPE` | Storage backend |
| `CLOTHO_METRICS_PORT` | Metrics server port |

Environment variables override config file values.

## Examples

### Workflow Definition (YAML)

```yaml
name: order-processing
version: "1.0"

description: Process customer orders with memory-augmented validation

activities:
  # Activity 1: Validate order
  - id: validate-order
    name: ValidateOrder
    timeout: 30s
    retry:
      max_attempts: 3
      initial_interval: 1s
      backoff: exponential

  # Activity 2: Check inventory (with memory lookup)
  - id: check-inventory
    name: CheckInventory
    timeout: 10s
    uses_memory: true  # Uses memory provider for context
    memory:
      provider: milvus
      query: "product_availability_{$.input.product_id}"

  # Activity 3: Process payment
  - id: process-payment
    name: ProcessPayment
    timeout: 60s
    retry:
      max_attempts: 5

  # Activity 4: Notify customer
  - id: notify-customer
    name: NotifyCustomer
    timeout: 15s

# Workflow DAG
dag:
  start: validate-order
  edges:
    - from: validate-order
      to: check-inventory
      condition: "$.valid == true"
    - from: check-inventory
      to: process-payment
    - from: process-payment
      to: notify-customer

# Compensation (Saga pattern)
compensation:
  - activity: process-payment
    handler: refund-payment
```

### Memory Provider Activity

```go
package activities

import (
    "context"
    "fmt"

    "github.com/Zhaoyikaiii/clotho/pkg/memory"
)

type CheckInventoryInput struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

type CheckInventoryOutput struct {
    Available bool    `json:"available"`
    Price     float64 `json:"price"`
    Context   string  `json:"context"` // From memory retrieval
}

func CheckInventory(ctx context.Context, input CheckInventoryInput) (*CheckInventoryOutput, error) {
    // Query memory layer for product context
    mem := memory.GetProvider("milvus")
    results, err := mem.Search(ctx, memory.SearchRequest{
        Collection: "product_knowledge",
        Query:      fmt.Sprintf("product %s availability stock", input.ProductID),
        TopK:       3,
    })
    if err != nil {
        return nil, fmt.Errorf("memory search failed: %w", err)
    }

    // Use memory context to enhance decision
    context := formatMemoryResults(results)

    // Check inventory (simplified)
    available := checkStock(input.ProductID, input.Quantity)

    return &CheckInventoryOutput{
        Available: available,
        Price:     getPrice(input.ProductID),
        Context:   context,
    }, nil
}
```

## Development

### Prerequisites

- Go 1.21+
- Docker & Docker Compose (for testing)

### Build

```bash
# Install dependencies
make deps

# Build binary
make build

# Run tests
make test

# Lint code
make lint
```

### Local Development

```bash
# Start local storage and memory services
docker compose -f docker/docker-compose.yml up -d

# Run the server
make run ARGS="server start"

# Run a workflow
make run ARGS="workflow run --name my-workflow"
```

### Project Structure

```
clotho/
├── cmd/
│   ├── clotho/              # Main CLI entry
│   ├── clotho-launcher/      # Launcher server
│   └── clotho-launcher-tui/ # TUI launcher
├── pkg/
│   ├── engine/              # Workflow orchestration engine
│   ├── worker/              # Worker implementation
│   ├── storage/             # State persistence
│   ├── memory/              # Memory/retrieval providers
│   ├── activities/          # Built-in activities
│   ├── observability/       # Tracing, metrics, logging
│   └── ...                  # Other packages
├── examples/                # Example workflows
├── docs/                    # Documentation
└── docker/                  # Docker configurations
```

## License

MIT License - see [LICENSE](./LICENSE) for details.

---

<div align="center">

Built with ❤️ by the Clotho team

</div>
