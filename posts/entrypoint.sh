#!/bin/sh
set -e
echo "Build app"

cd /app
go mod download
go build -o /app/go_app
echo "Running app"
exec /app/go_app
