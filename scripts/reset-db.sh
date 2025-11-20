#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${RED}========================================${NC}"
echo -e "${RED}WARNING: This will delete all data!${NC}"
echo -e "${RED}========================================${NC}"
echo ""
echo -e "${YELLOW}This script will:${NC}"
echo "1. Drop all migrations"
echo "2. Delete the database"
echo "3. Recreate the database"
echo "4. Run all migrations from scratch"
echo ""
read -p "Are you sure? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "Cancelled."
    exit 0
fi

# Drop all migrations
echo -e "${YELLOW}[1/4]${NC} Dropping all migrations..."
migrate -path ./migrations -database "postgres://healthuser:healthpass@localhost:5432/healthdb?sslmode=disable" drop -f
echo -e "${GREEN}✓ Migrations dropped${NC}"

# Drop database
echo -e "${YELLOW}[2/4]${NC} Dropping database..."
sudo -u postgres psql << EOF
DROP DATABASE IF EXISTS healthdb;
\q
EOF
echo -e "${GREEN}✓ Database dropped${NC}"

# Create database
echo -e "${YELLOW}[3/4]${NC} Creating database..."
sudo -u postgres psql << EOF
CREATE DATABASE healthdb OWNER healthuser;
GRANT ALL PRIVILEGES ON DATABASE healthdb TO healthuser;
\q
EOF
echo -e "${GREEN}✓ Database created${NC}"

# Run migrations
echo -e "${YELLOW}[4/4]${NC} Running migrations..."
migrate -path ./migrations -database "postgres://healthuser:healthpass@localhost:5432/healthdb?sslmode=disable" up
echo -e "${GREEN}✓ Migrations completed${NC}"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Database reset complete!${NC}"
echo -e "${GREEN}========================================${NC}"