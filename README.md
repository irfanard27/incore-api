# Incore API

Inventory management API built with Go and PostgreSQL.

## Prerequisites

- Go 1.25.0 or higher
- PostgreSQL database
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd incore-api
```

2. Install dependencies:
```bash
go mod download
```

3. Set up the database:
```bash
# Create database
createdb incore-db

# Run migrations (if migration files exist in db/migrations/)
```

4. Copy environment file:
```bash
cp .env.example .env
```

5. Update `.env` file with your database configuration:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=incore-db
JWT_SECRET=your_jwt_secret
PORT=8080
```

## Running the Application

### Development Mode

```bash
go run cmd/server/main.go
```

### Production Mode

```bash
# Build the application
go build -o incore-api cmd/server/main.go

# Run the binary
./incore-api
```

## API Endpoints

The API will be available at `http://localhost:8080`

### Authentication
- `POST /auth/login` - User login
- `POST /auth/register` - User registration

### Inventory Management
- `GET /inventory` - Get all inventories
- `GET /inventory/:id` - Get inventory by ID
- `POST /inventory` - Create new inventory
- `PUT /inventory` - Update inventory
- `DELETE /inventory/:id` - Delete inventory

### Stock Management
- `POST /stock-out` - Create stock out
- `POST /stock-out/:id/rollback` - Rollback stock out
- `POST /stock-in` - Create stock in

## Database Schema

The application uses PostgreSQL with the following main tables:
- `users` - User accounts
- `inventories` - Inventory items
- `stock_out` - Stock out transactions
- `stock_out_items` - Stock out item details
- `stock_in` - Stock in transactions

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database username | postgres |
| DB_PASSWORD | Database password | postgres |
| DB_NAME | Database name | incore-db |
| JWT_SECRET | JWT secret key | - |
| PORT | Server port | 8080 |

## Project Structure

```
incore-api/
├── cmd/server/          # Main application entry point
├── internal/
│   ├── delivery/http/   # HTTP handlers and routes
│   ├── domain/          # Domain entities and DTOs
│   ├── infra/           # Infrastructure (database, JWT)
│   └── usecase/         # Business logic
├── db/migrations/       # Database migration files
├── pkg/                 # Shared packages
└── .env                 # Environment variables
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test your changes
5. Submit a pull request

## License

[Add your license information here]
