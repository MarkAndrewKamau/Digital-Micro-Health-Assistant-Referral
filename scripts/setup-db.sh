#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}Digital Health Assistant - Database Setup${NC}"
echo -e "${YELLOW}========================================${NC}"

# Check if PostgreSQL is running
echo -e "${YELLOW}[1/6]${NC} Checking PostgreSQL status..."
if ! sudo systemctl is-active --quiet postgresql; then
    echo -e "${YELLOW}Starting PostgreSQL...${NC}"
    sudo systemctl start postgresql
    sleep 2
else
    echo -e "${GREEN}✓ PostgreSQL is already running${NC}"
fi

# Check if migrate tool is installed
echo -e "${YELLOW}[2/6]${NC} Checking migrate tool..."
if ! command -v migrate &> /dev/null; then
    echo -e "${RED}✗ migrate tool not found${NC}"
    echo -e "${YELLOW}Installing migrate...${NC}"
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
    sudo mv migrate /usr/local/bin/migrate
    echo -e "${GREEN}✓ migrate installed${NC}"
else
    echo -e "${GREEN}✓ migrate is already installed${NC}"
fi

# Drop and recreate database
echo -e "${YELLOW}[3/6]${NC} Dropping existing database and user..."
sudo -u postgres psql << EOF
DROP DATABASE IF EXISTS healthdb;
DROP USER IF EXISTS healthuser;
\q
EOF
echo -e "${GREEN}✓ Dropped existing database${NC}"

# Create new user and database
echo -e "${YELLOW}[4/6]${NC} Creating new database and user..."
sudo -u postgres psql << EOF
CREATE USER healthuser WITH PASSWORD 'healthpass';
CREATE DATABASE healthdb OWNER healthuser;
GRANT ALL PRIVILEGES ON DATABASE healthdb TO healthuser;
ALTER USER healthuser WITH SUPERUSER;
\q
EOF
echo -e "${GREEN}✓ Database and user created${NC}"

# Run migrations
echo -e "${YELLOW}[5/6]${NC} Running migrations..."
migrate -path ./migrations -database "postgres://healthuser:healthpass@localhost:5432/healthdb?sslmode=disable" up
echo -e "${GREEN}✓ Migrations completed${NC}"

# Verify tables
echo -e "${YELLOW}[6/6]${NC} Verifying database tables..."
psql -U healthuser -d healthdb -h localhost -c "\dt"
echo -e "${GREEN}✓ Database setup complete${NC}"

# Start Redis
echo -e "${YELLOW}Starting Redis...${NC}"
sudo systemctl start redis-server
sleep 1

if redis-cli ping | grep -q "PONG"; then
    echo -e "${GREEN}✓ Redis is running${NC}"
else
    echo -e "${RED}✗ Redis failed to start${NC}"
    exit 1
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Database setup complete!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Database credentials:${NC}"
echo "  User: healthuser"
echo "  Password: healthpass"
echo "  Database: healthdb"
echo "  Host: localhost:5432"
echo ""
echo -e "${YELLOW}Migration commands:${NC}"
echo "  Up:      migrate -path ./migrations -database \"postgres://healthuser:healthpass@localhost:5432/healthdb?sslmode=disable\" up"
echo "  Down:    migrate -path ./migrations -database \"postgres://healthuser:healthpass@localhost:5432/healthdb?sslmode=disable\" down"
echo "  Version: migrate -path ./migrations -database \"postgres://healthuser:healthpass@localhost:5432/healthdb?sslmode=disable\" version"