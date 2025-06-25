# Yoyo Store Backend

Yoyo Store provides the backend services for a sample e‑commerce shop selling yoyos. The project is written in Go and exposes a REST API together with a gRPC microservice for generating invoices.
It exposes REST APIs for the storefront and a gRPC service for invoice generation. Payments and subscriptions are processed with Stripe while all components run in containers and are deployed via GitHub Actions to Amazon EKS.
## Features

- RESTful API in Golang
- **Stripe integration** – charge single items, subscribe customers to recurring plans, refund payments and cancel subscriptions.
- **Invoice service** – a dedicated gRPC service creates PDF invoices and emails them to customers.
- **Email notifications** – send payment receipts and forgotten password links via SMTP.
- **GitHub Actions CI/CD** – linting, tests and Docker builds run automatically. Images are deployed to Amazon EKS and exposed with AWS Application Load Balancers.
- **AWS load balancers** – production deployments use Application Load Balancer for reliability.
- User authentication (JWT)
- Admin endpoints

## Demo

A short demo is available in [`yoyo_store.mp4`](yoyo_store.mp4).

<video src="yoyo_store.mp4" controls width="600"></video>

## Tech Stack

- **Language:** Golang
- **Framework:** Chi router
- **Database:** PostgreSQL
- **Payments:** Stripe
- **Messaging:** gRPC for invoice generation
- **Containerization:** Docker
- **Orchestration:** Kubernetes on Amazon EKS

## Prerequisites

- [Go](https://golang.org/doc/install) ≥ 1.18
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

Create a `.env` file in the project root and populate it with the following values:

```env
ENVIRONMENT=develop
ALLOWED_ORIGINS=http://localhost:3000
DB_SOURCE=postgresql://root:secret@localhost:5432/yoyo_store?sslmode=disable
MIGRATION_URL=file://db/migration
MAIN_SERVER_PORT=8080
INVOICE_GRPC_PORT=9090
INVOICE_HTTP_PORT=9091
FRONTEND_PORT=3000
TOKEN_SYMMETRIC_KEY=your-secret-key
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USERNAME=
SMTP_PASSWORD=
MAIN_SERVER_ADDR=http://localhost:8080
INVOICE_GRPC_ADDR=http://localhost:9090
FRONTEND_ADDR=http://localhost:3000
STRIPE_SECRET=
STRIPE_KEY=
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

### Testing

- **Run tests:**
```bash
make test
```

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

## Docker

Each service has its own `Dockerfile` for building a container image.

### Build images

```bash
make build_docker_back
make build_docker_invoice
make build_docker_front
make build_docker
```

Run the containers with your environment variables and expose the required ports:

```bash
make run_docker_back
make run_docker_invoice
make run_docker_front
```

- **Docker Compose:**
Run all services with Docker Compose:
```bash
make run-compose-local
make stop-compose-local
```

- **Useful Docker commands:**
```bash
docker ps
docker rm {container-name}
docker rmi {image-id}
docker container inspect {container-name}
docker network create {network-name}
docker network connect {network-name} {container-name}
docker network ls
docker network inspect {network-name}
docker stop $(docker ps -a -q)
docker rm -f $(docker ps -a -q)
docker rmi -f $(docker images -aq)
```

## CI/CD

- `test.yaml` runs linting and unit tests on every push and pull request.
- `deploy-staging.yaml` builds Docker images, pushes them to Amazon ECR and deploys to an EKS cluster using AWS load balancers.

## gRPC Invoice Service

The invoice microservice exposes a `CreateAndSendInvoice` RPC defined in [`internal/proto/invoice.proto`](internal/proto/invoice.proto). After a successful payment, the main server calls this endpoint to generate a PDF invoice and send it by email.

## Stripe Integration

The API supports:

- Obtaining payment intents to charge items.
- Subscribing and unsubscribing customers from plans.
- Refunding charges.

## Email Notifications

Emails are delivered through SMTP for purchase receipts and password reset requests.

## License

This project is licensed under the terms of the MIT license.