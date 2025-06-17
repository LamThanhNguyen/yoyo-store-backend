# Yoyo Store Backend

Backend system for a yoyo e-commerce store, built with Golang. 
Features user management, product catalog, order processing, Stripe-powered payments and subscriptions, invoicing, and refunds.

## Features

- RESTful API in Golang
- Stripe integration for payments & subscriptions
- Invoice generation and management
- Refund processing
- User authentication (JWT)
- Product catalog & inventory
- Order management
- Subscription plans (recurring payments)
- Admin endpoints

## Tech Stack

- **Language:** Golang
- **Framework:** Gin / Fiber / Echo (TBD)
- **Database:** PostgreSQL / MongoDB (TBD)
- **Payments:** Stripe
- **Deployment:** Docker, Kubernetes (optional)

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.18
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/)
- [PostgreSQL](https://www.postgresql.org/)

### Installation

#### Install Required Tools

- **Migrate** ([docs](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)):
    ```bash
    curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | sudo apt-key add -
    echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" | sudo tee /etc/apt/sources.list.d/migrate.list
    sudo apt-get update
    sudo apt-get install -y migrate
    ```

- **GoMock** ([docs](https://github.com/uber-go/mock)):
    ```bash
    go install go.uber.org/mock/mockgen@latest
    export PATH=$PATH:$(go env GOPATH)/bin
    mockgen -version
    ```

### Environment Variables

Create a `.env` file in the project root and fill in the following:

```env
ENVIRONMENT=develop
ALLOWED_ORIGINS=http://localhost:3000
DB_SOURCE=postgresql://{{username}}:{{password}}@localhost:5432/{{database_name}}?sslmode=disable
MIGRATION_URL=file://db/migration
REDIS_ADDRESS=localhost:6379
HTTP_SERVER_ADDRESS=0.0.0.0:8080
TOKEN_SYMMETRIC_KEY=2e3c226355a0770689c808684fbdca40
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=24h
EMAIL_SENDER_NAME=
EMAIL_SENDER_ADDRESS=
EMAIL_SENDER_PASSWORD=
FRONTEND_DOMAIN=http://localhost:3000
```

### Database & Infrastructure

- **Create Docker network:**
    ```bash
    make network
    ```

- **Start PostgreSQL:**
    ```bash
    make postgres
    ```

- **Create database:**
    ```bash
    make createdb
    ```

- **Run migrations:**
    ```bash
    make migrateup      # Up all versions
    make migrateup1     # Up 1 version
    make migratedown    # Down all versions
    make migratedown1   # Down 1 version
    ```

### Code Generation

- **Create a new DB migration:**
    ```bash
    make new_migration name=<migration_name>
    ```

- **Initialize Go module:**
    ```bash
    go mod init github.com/LamThanhNguyen/yoyo-store-backend
    ```

- **Install Go packages:**
    ```bash
    go get github.com/some/library
    go mod tidy
    ```

- **Generate DB mocks with GoMock:**
    ```bash
    make mock
    ```

### Running the Application
- **Ensure you already install swag:**
    ```bash
    go install github.com/swaggo/swag/cmd/swag@latest
    echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc
    source ~/.bashrc
    swag --version
    // go get -u github.com/swaggo/gin-swagger
    // go get -u github.com/swaggo/files
    ```

- **Run server:**
    ```bash
    make server
    ```

### Testing

- **Run tests:**
    ```bash
    make test
    ```

## API Documentation

- **Generate Swagger docs:**
    ```bash
    swag init -g main.go --output docs
    ```
- **View docs:**  
  Visit [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) after running the server.

## Linting

- **Install golangci-lint:**
    ```bash
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.1.6
    echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc
    source ~/.bashrc
    ```

- **Run linter:**
    ```bash
    golangci-lint --version
    golangci-lint run
    ```