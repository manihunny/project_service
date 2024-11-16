#!/bin/sh

echo "Waiting for database to be ready..."
wait-for "db:5432" -- "$@"

echo "Waiting for Redis to be ready..."
wait-for "redis:6379" -- "$@"

echo "Running database migrations..."
go run ./cmd/migrations/main.go

echo "Starting the application..."
CompileDaemon --build="go build -mod vendor -o bin/main ./cmd/service/main.go"  --command=./bin/main