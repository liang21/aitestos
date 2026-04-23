#!/bin/bash
# Migration runner for aitestos database

set -e

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS_DIR="$SCRIPT_DIR/migrations"

# Default database connection string (can be overridden by env vars)
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_NAME="${DB_NAME:-aitestos}"
DB_USER="${DB_USER:-postgres}"
DB_PASSWORD="${DB_PASSWORD:-postgres}"

# Build connection string
DB_CONN="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"

# Command: up or down
COMMAND=${1:-up}

case $COMMAND in
    up)
        echo "Running database migrations..."
        for migration in "$MIGRATIONS_DIR"/*.sql; do
            if [ -f "$migration" ]; then
                echo "Applying: $(basename "$migration")"
                psql "$DB_CONN" -f "$migration"
            fi
        done
        echo "Migrations completed!"
        ;;
    down)
        echo "Rolling back migrations is not implemented yet."
        echo "Please manually execute SQL rollback statements."
        exit 1
        ;;
    status)
        echo "Checking migration status..."
        psql "$DB_CONN" -c "\dt"
        ;;
    *)
        echo "Usage: $0 {up|down|status}"
        exit 1
        ;;
esac
