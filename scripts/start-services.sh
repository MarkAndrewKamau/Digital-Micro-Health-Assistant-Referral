#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}Starting Services${NC}"
echo -e "${YELLOW}========================================${NC}"

# Start PostgreSQL
echo -e "${YELLOW}[1/2]${NC} Starting PostgreSQL..."
sudo systemctl start postgresql
sleep 2

if sudo systemctl is-active --quiet postgresql; then
    echo -e "${GREEN}✓ PostgreSQL started${NC}"
else
    echo -e "${RED}✗ PostgreSQL failed to start${NC}"
    exit 1
fi

# Start Redis
echo -e "${YELLOW}[2/2]${NC} Starting Redis..."
sudo systemctl start redis-server
sleep 1

if redis-cli ping | grep -q "PONG"; then
    echo -e "${GREEN}✓ Redis started${NC}"
else
    echo -e "${RED}✗ Redis failed to start${NC}"
    exit 1
fi

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}All services started successfully!${NC}"
echo -e "${GREEN}========================================${NC}"