#!/bin/bash

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

DATABASE_URL=${DATABASE_URL:-"postgres://healthuser:healthpass@localhost:5432/healthdb?sslmode=disable"}
MIGRATIONS_PATH="./migrations"

case "$1" in
    up)
        echo "Running migrations up..."
        migrate -path $MIGRATIONS_PATH -database "$DATABASE_URL" up
        ;;
    down)
        echo "Running migrations down..."
        migrate -path $MIGRATIONS_PATH -database "$DATABASE_URL" down
        ;;
    drop)
        echo "Dropping all migrations..."
        migrate -path $MIGRATIONS_PATH -database "$DATABASE_URL" drop -f
        ;;
    force)
        if [ -z "$2" ]; then
            echo "Please provide version number: ./migrate.sh force VERSION"
            exit 1
        fi
        echo "Forcing migration version to $2..."
        migrate -path $MIGRATIONS_PATH -database "$DATABASE_URL" force $2
        ;;
    version)
        echo "Current migration version:"
        migrate -path $MIGRATIONS_PATH -database "$DATABASE_URL" version
        ;;
    create)
        if [ -z "$2" ]; then
            echo "Please provide migration name: ./migrate.sh create NAME"
            exit 1
        fi
        echo "Creating migration: $2"
        migrate create -ext sql -dir $MIGRATIONS_PATH -seq $2
        ;;
    *)
        echo "Usage: $0 {up|down|drop|force|version|create}"
        exit 1
        ;;
esac