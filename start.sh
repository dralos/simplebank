#!/bin/sh
set -e

echo "running database migration..."
/app/migrate -path /app/migration -database "${DB_SOURCE}" -verbose up

echo "starting the application..."
exec "$@"