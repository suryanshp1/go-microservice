# 🚀 Go Microservice E-Commerce Platform

A production-grade microservice-based e-commerce platform built with **Go**, featuring **gRPC** inter-service communication, **GraphQL** API gateway, **PostgreSQL** & **Elasticsearch** databases, all orchestrated with **Docker Compose**.

---

## 📋 Table of Contents

- [Executive Summary](#executive-summary)
- [Architecture Overview](#architecture-overview)
- [Tech Stack](#tech-stack)
- [Service Breakdown](#service-breakdown)
- [Getting Started](#getting-started)
- [API Reference](#api-reference)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Development Guide](#development-guide)
- [Risks & Considerations](#risks--considerations)

---

## Executive Summary

This project demonstrates a **cloud-native microservices architecture** implementing an e-commerce domain with:

| Component | Technology |
|-----------|------------|
| **API Gateway** | GraphQL (gqlgen) |
| **Inter-Service Communication** | gRPC with Protocol Buffers |
| **Account Database** | PostgreSQL |
| **Catalog Database** | Elasticsearch |
| **Order Database** | PostgreSQL |
| **Containerization** | Docker Compose |

The system follows **Domain-Driven Design (DDD)** principles with clear service boundaries, ensuring scalability, maintainability, and independent deployment capabilities.

---

## Architecture Overview

### High-Level System Architecture

```mermaid
flowchart TB
    subgraph Client Layer
        WEB[🌐 Web Client]
        MOBILE[📱 Mobile Client]
    end
    
    subgraph API Gateway
        GQL[📡 GraphQL Gateway<br/>Port: 8000]
    end
    
    subgraph Microservices
        ACC[👤 Account Service<br/>gRPC :8080]
        CAT[📦 Catalog Service<br/>gRPC :8080]
        ORD[🛒 Order Service<br/>gRPC :8080]
    end
    
    subgraph Data Layer
        PG_ACC[(🐘 PostgreSQL<br/>Account DB)]
        ES[(🔍 Elasticsearch<br/>Catalog DB)]
        PG_ORD[(🐘 PostgreSQL<br/>Order DB)]
    end
    
    WEB --> GQL
    MOBILE --> GQL
    
    GQL -->|gRPC| ACC
    GQL -->|gRPC| CAT
    GQL -->|gRPC| ORD
    
    ORD -->|gRPC| ACC
    ORD -->|gRPC| CAT
    
    ACC --> PG_ACC
    CAT --> ES
    ORD --> PG_ORD
    
    classDef gateway fill:#e1f5fe,stroke:#01579b
    classDef service fill:#f3e5f5,stroke:#4a148c
    classDef database fill:#e8f5e9,stroke:#1b5e20
    
    class GQL gateway
    class ACC,CAT,ORD service
    class PG_ACC,ES,PG_ORD database
```

### Service Communication Flow

```mermaid
sequenceDiagram
    autonumber
    participant C as Client
    participant G as GraphQL Gateway
    participant A as Account Service
    participant P as Catalog Service
    participant O as Order Service
    participant DB as Databases

    rect rgb(240, 248, 255)
        Note over C,G: Query Flow
        C->>G: GraphQL Query
        G->>A: GetAccount (gRPC)
        A->>DB: SELECT FROM accounts
        DB-->>A: Account Data
        A-->>G: Account Response
        G-->>C: GraphQL Response
    end

    rect rgb(255, 248, 240)
        Note over C,O: Order Creation Flow
        C->>G: createOrder Mutation
        G->>O: PostOrder (gRPC)
        O->>A: Validate Account
        A-->>O: Account Valid
        O->>P: Get Product Details
        P-->>O: Product Data
        O->>DB: INSERT INTO orders
        DB-->>O: Order Created
        O-->>G: Order Response
        G-->>C: GraphQL Response
    end
```

### Container Architecture

```mermaid
flowchart LR
    subgraph Docker Network
        subgraph Services
            GQL[graphql<br/>:8000→:8080]
            ACC[account<br/>:8080]
            CAT[catalog<br/>:8080]
            ORD[order<br/>:8080]
        end
        
        subgraph Databases
            ACC_DB[(account_db<br/>PostgreSQL)]
            CAT_DB[(catalog_db<br/>Elasticsearch)]
            ORD_DB[(order_db<br/>PostgreSQL)]
        end
    end
    
    EXT[External<br/>:8000] --> GQL
    
    GQL -.-> ACC
    GQL -.-> CAT
    GQL -.-> ORD
    
    ORD -.-> ACC
    ORD -.-> CAT
    
    ACC --> ACC_DB
    CAT --> CAT_DB
    ORD --> ORD_DB
```

### Data Models

```mermaid
erDiagram
    ACCOUNT {
        string id PK
        string name
    }
    
    PRODUCT {
        string id PK
        string name
        string description
        float price
    }
    
    ORDER {
        string id PK
        string account_id FK
        timestamp created_at
        timestamp updated_at
        float total_price
    }
    
    ORDER_PRODUCT {
        string order_id FK
        string product_id FK
        int quantity
    }
    
    ACCOUNT ||--o{ ORDER : places
    ORDER ||--o{ ORDER_PRODUCT : contains
    PRODUCT ||--o{ ORDER_PRODUCT : included_in
```

---

## Tech Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| **Language** | Go 1.24.4 | Core service implementation |
| **API Gateway** | gqlgen v0.17.85 | GraphQL schema-first development |
| **RPC Framework** | gRPC + Protobuf | High-performance service-to-service communication |
| **Account DB** | PostgreSQL | ACID-compliant relational data storage |
| **Catalog DB** | Elasticsearch 6.2.4 | Full-text search & product catalog |
| **Order DB** | PostgreSQL | Transactional order processing |
| **Container** | Docker Compose | Multi-container orchestration |
| **ID Generation** | KSUID | Sortable, unique identifiers |

---

## Service Breakdown

### 🌐 GraphQL Gateway (`/graphql`)

The **API Gateway** exposes a unified GraphQL endpoint for all client interactions.

**Key Features:**
- Schema-first GraphQL design
- Playground available at `/playground`
- Aggregates data from all microservices
- Connection pooling to backend services

**Exposed Endpoints:**
| Endpoint | Description |
|----------|-------------|
| `POST /graphql` | GraphQL API endpoint |
| `GET /playground` | GraphQL Playground IDE |

**GraphQL Schema:**
```graphql
type Query {
  accounts(pagination: PaginationInput, id: String): [Account!]!
  products(pagination: PaginationInput, query: String, id: String): [Product!]!
}

type Mutation {
  createAccount(input: AccountInput!): Account
  createProduct(input: ProductInput!): Product
  createOrder(input: OrderInput!): Order
}
```

---

### 👤 Account Service (`/account`)

Manages user account lifecycle with **PostgreSQL** persistence.

**Architecture Pattern:** Repository Pattern with Service Layer

```mermaid
flowchart LR
    CLIENT[gRPC Client] --> SERVER[gRPC Server]
    SERVER --> SERVICE[Account Service]
    SERVICE --> REPO[PostgreSQL Repository]
    REPO --> DB[(PostgreSQL)]
```

**API Operations:**
| RPC Method | Description |
|------------|-------------|
| `PostAccount` | Create new account |
| `GetAccount` | Retrieve account by ID |
| `GetAccounts` | List accounts with pagination |

**Data Model:**
```go
type Account struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

---

### 📦 Catalog Service (`/catalog`)

Product catalog management powered by **Elasticsearch** for full-text search.

**Key Features:**
- Full-text product search
- Batch product retrieval by IDs
- Pagination support

**API Operations:**
| RPC Method | Description |
|------------|-------------|
| `PostProduct` | Add new product |
| `GetProduct` | Get product by ID |
| `GetProducts` | Search/list products with filters |

**Data Model:**
```go
type Product struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
}
```

---

### 🛒 Order Service (`/order`)

Orchestrates order processing with **cross-service coordination**.

**Key Features:**
- Account validation via Account Service
- Product enrichment via Catalog Service
- Order persistence in PostgreSQL
- Total price calculation

**API Operations:**
| RPC Method | Description |
|------------|-------------|
| `PostOrder` | Create new order |
| `GetOrdersForAccount` | Get all orders for an account |

**Order Flow:**
```mermaid
flowchart TD
    A[Receive Order Request] --> B{Validate Account}
    B -->|Invalid| C[Return Error]
    B -->|Valid| D[Fetch Product Details]
    D --> E[Calculate Total Price]
    E --> F[Persist Order]
    F --> G[Return Order Response]
```

---

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Go 1.24+ (for development)
- Protocol Buffer Compiler (for proto changes)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/suryanshp1/go-microservice.git
cd go-microservice

# Start all services
docker-compose up --build

# Access GraphQL Playground
open http://localhost:8000/playground
```

### Service Startup Order

```mermaid
flowchart TD
    A[account_db] --> B[account]
    C[catalog_db] --> D[catalog]
    E[order_db] --> F[order]
    B --> G[graphql]
    D --> G
    F --> G
```

---

## API Reference

### GraphQL Examples

#### Create Account
```graphql
mutation {
  createAccount(input: { name: "John Doe" }) {
    id
    name
  }
}
```

#### Create Product
```graphql
mutation {
  createProduct(input: {
    name: "Laptop"
    description: "High-performance laptop"
    price: 999.99
  }) {
    id
    name
    price
  }
}
```

#### Create Order
```graphql
mutation {
  createOrder(input: {
    accountId: "2KxQjP7mGqSz8tB1n4v3wR"
    products: [
      { id: "product-id-1", quantity: 2 }
      { id: "product-id-2", quantity: 1 }
    ]
  }) {
    id
    totalPrice
    products {
      name
      quantity
      price
    }
  }
}
```

#### Query Accounts with Orders
```graphql
query {
  accounts(id: "account-id") {
    id
    name
    orders {
      id
      totalPrice
      createdAt
      products {
        name
        price
        quantity
      }
    }
  }
}
```

#### Search Products
```graphql
query {
  products(query: "laptop", pagination: { skip: 0, take: 10 }) {
    id
    name
    description
    price
  }
}
```

---

## Project Structure

```
go-microservice/
├── 📁 account/              # Account Microservice
│   ├── account.proto        # gRPC service definition
│   ├── app.dockerfile       # Service container
│   ├── db.dockerfile        # PostgreSQL container
│   ├── client.go            # gRPC client
│   ├── server.go            # gRPC server
│   ├── service.go           # Business logic
│   ├── repository.go        # Data access layer
│   ├── up.sql               # Database migrations
│   ├── cmd/                  # Service entrypoint
│   └── pb/                   # Generated protobuf code
│
├── 📁 catalog/              # Catalog Microservice
│   ├── catalog.proto        # gRPC service definition
│   ├── app.dockerfile       # Service container
│   ├── client.go            # gRPC client
│   ├── server.go            # gRPC server
│   ├── service.go           # Business logic
│   ├── repository.go        # Elasticsearch repository
│   ├── cmd/                  # Service entrypoint
│   └── pb/                   # Generated protobuf code
│
├── 📁 order/                # Order Microservice
│   ├── order.proto          # gRPC service definition
│   ├── app.dockerfile       # Service container
│   ├── db.dockerfile        # PostgreSQL container
│   ├── client.go            # gRPC client
│   ├── server.go            # gRPC server + orchestration
│   ├── service.go           # Business logic
│   ├── repository.go        # Data access layer
│   ├── up.sql               # Database migrations
│   ├── cmd/                  # Service entrypoint
│   └── pb/                   # Generated protobuf code
│
├── 📁 graphql/              # GraphQL Gateway
│   ├── schema.graphql       # GraphQL schema
│   ├── main.go              # HTTP server entrypoint
│   ├── graph.go             # GraphQL server setup
│   ├── *_resolver.go        # GraphQL resolvers
│   ├── models_gen.go        # Generated models
│   ├── generated.go         # Generated runtime
│   └── gqlgen.yml           # Code generation config
│
├── 📁 diagrams/             # Architecture diagrams
├── 📁 vendor/               # Go dependencies
├── docker-compose.yaml      # Container orchestration
├── go.mod                   # Go module definition
└── go.sum                   # Dependency checksums
```

---

## Configuration

### Environment Variables

| Service | Variable | Description | Default |
|---------|----------|-------------|---------|
| **account** | `DATABASE_URL` | PostgreSQL connection string | - |
| **catalog** | `DATABASE_URL` | Elasticsearch URL | - |
| **order** | `DATABASE_URL` | PostgreSQL connection string | - |
| **order** | `ACCOUNT_SERVICE_URL` | Account service gRPC address | - |
| **order** | `CATALOG_SERVICE_URL` | Catalog service gRPC address | - |
| **graphql** | `ACCOUNT_SERVICE_URL` | Account service gRPC address | - |
| **graphql** | `CATALOG_SERVICE_URL` | Catalog service gRPC address | - |
| **graphql** | `ORDER_SERVICE_URL` | Order service gRPC address | - |

---

## Development Guide

### Modifying Protobuf Definitions

```bash
# Install protoc plugins
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Regenerate Go code (example for account service)
cd account
protoc --go_out=./ --go-grpc_out=./ account.proto
```

### Regenerating GraphQL Code

```bash
cd graphql
go run github.com/99designs/gqlgen generate
```

### Running Individual Services

```bash
# Account Service
DATABASE_URL=postgres://user:pass@localhost/db go run ./account/cmd

# Catalog Service
DATABASE_URL=http://localhost:9200 go run ./catalog/cmd

# Order Service
DATABASE_URL=postgres://user:pass@localhost/db \
ACCOUNT_SERVICE_URL=localhost:8080 \
CATALOG_SERVICE_URL=localhost:8081 \
go run ./order/cmd
```

---

## Risks & Considerations

| Risk | Impact | Mitigation |
|------|--------|------------|
| **Single Point of Failure (GraphQL Gateway)** | High | Deploy multiple instances with load balancer |
| **No Authentication/Authorization** | Critical | Implement JWT/OAuth2 before production |
| **Hardcoded Credentials** | Critical | Use secrets management (Vault, K8s secrets) |
| **No Rate Limiting** | Medium | Add rate limiting at gateway level |
| **Elasticsearch v6.2.4 EOL** | Medium | Upgrade to Elasticsearch 8.x |
| **No Distributed Tracing** | Medium | Integrate OpenTelemetry/Jaeger |
| **No Health Checks** | Medium | Implement /health endpoints |
| **No Circuit Breaker** | Medium | Add resilience patterns (hystrix, go-kit) |

---

## 📄 License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
