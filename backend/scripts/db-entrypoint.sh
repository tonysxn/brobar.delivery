#!/bin/sh
set -ex

DB_HOST=${DB_HOST:-db}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-${SERVICE_NAME}_db}
MIGRATIONS_PATH=${MIGRATIONS_PATH:-/app/migrations}

until pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER"; do
  echo "Waiting for PostgreSQL at ${DB_HOST}:${DB_PORT}..."
  sleep 2
done

PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<-EOSQL
    CREATE TABLE IF NOT EXISTS schema_migrations (
        version bigint NOT NULL PRIMARY KEY,
        dirty boolean NOT NULL
    );
EOSQL

if [ -d "$MIGRATIONS_PATH" ]; then
  echo "Applying migrations from ${MIGRATIONS_PATH}"
  migrate -path="$MIGRATIONS_PATH" -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up
fi

exec "$@"