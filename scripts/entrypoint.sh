#!/bin/sh

echo "Waiting for database to be ready..."
wait-for "db:5432" -- "$@"

echo "Running database migrations..."
./bin/migrations

echo "Starting the application..."
./bin/main