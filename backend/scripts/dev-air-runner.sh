#!/bin/sh
# Usage: ./scripts/dev-air-runner.sh <SERVICE_NAME> <MAIN_PATH>
# Example: ./scripts/dev-air-runner.sh user-service cmd/user/main.go

SERVICE_NAME=$1
MAIN_PATH=$2

if [ -z "$SERVICE_NAME" ] || [ -z "$MAIN_PATH" ]; then
    echo "Error: SERVICE_NAME and MAIN_PATH arguments are required."
    exit 1
fi

CONFIG_PATH="/tmp/air-${SERVICE_NAME}.toml"

echo "Generating air config for $SERVICE_NAME at $CONFIG_PATH..."

sed -e "s|__SERVICE_NAME__|$SERVICE_NAME|g" \
    -e "s|__MAIN_PATH__|$MAIN_PATH|g" \
    air.toml.tmpl > "$CONFIG_PATH"

echo "Starting air..."
air -c "$CONFIG_PATH"
