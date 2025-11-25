# Development Setup

## First Time Setup (Fresh Install)

Run this once to set up everything:
```bash
./scripts/setup-db.sh
go mod tidy
go run cmd/server/main.go
```

## After Reboot (Live Boot Environment)

Since you're on a live boot, run this after each reboot to start services:
```bash
./scripts/start-services.sh
go mod tidy
go run cmd/server/main.go
```

## Database Commands
```bash
# View database tables
psql -U healthuser -d healthdb -h localhost -c "\dt"

# View specific table
psql -U healthuser -d healthdb -h localhost -c "SELECT * FROM users;"

# Connect to database shell
psql -U healthuser -d healthdb -h localhost

# Reset database (WARNING: deletes all data)
./scripts/reset-db.sh

# Check migrations version
migrate -path ./migrations -database "postgres://healthuser:healthpass@localhost:5432/healthdb?sslmode=disable" version
```

## Database Credentials

- **User**: healthuser
- **Password**: healthpass
- **Database**: healthdb
- **Host**: localhost
- **Port**: 5432

## Services

Check status:
```bash
sudo systemctl status postgresql
sudo systemctl status redis-server
```

Start manually:
```bash
sudo systemctl start postgresql
sudo systemctl start redis-server
```

Stop:
```bash
sudo systemctl stop postgresql
sudo systemctl stop redis-server
```

## Testing
```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -v -run TestLoginFailed ./internal/handlers
```

## Common Issues

### "password authentication failed"
- Ensure PostgreSQL is running: `sudo systemctl start postgresql`
- Check credentials in `.env` file

### "connection refused" (Redis)
- Ensure Redis is running: `sudo systemctl start redis-server`
- Test with: `redis-cli ping`

### Migrations failed
- Reset database: `./scripts/reset-db.sh`
- Or manually: `migrate -path ./migrations -database "postgres://..." down` then `up`

## API Testing
```bash
# Register a user
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"phone": "+254712345678", "role": "patient"}'

# Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"phone": "+254712345678"}'

# Use session token
curl http://localhost:8080/v1/auth/me \
  -H "Authorization: Bearer YOUR_SESSION_TOKEN"
```