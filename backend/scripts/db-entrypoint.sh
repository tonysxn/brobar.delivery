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


if [ -d "$MIGRATIONS_PATH" ]; then
  echo "Applying migrations from ${MIGRATIONS_PATH}"
  CONNECTION="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
  if ! migrate -path="$MIGRATIONS_PATH" -database "$CONNECTION" up; then
    echo "Migration failed. Checking for dirty state..."
    DIRTY_VERSION=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT version FROM schema_migrations WHERE dirty = true LIMIT 1" | xargs)
    if [ -n "$DIRTY_VERSION" ]; then
      echo "Database is dirty at version $DIRTY_VERSION. Forcing and retrying..."
      migrate -path="$MIGRATIONS_PATH" -database "$CONNECTION" force "$DIRTY_VERSION"
      migrate -path="$MIGRATIONS_PATH" -database "$CONNECTION" up
    else
      echo "Migration failed with a non-dirty error. Check logs above."
      exit 1
    fi
  fi
fi

exec "$@"