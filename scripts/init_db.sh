#!/bin/bash

# Create database
createdb llama3_db

# Create user if not exists
psql -c "CREATE USER postgres WITH PASSWORD 'postgres' SUPERUSER;" postgres

# Grant privileges
psql -c "GRANT ALL PRIVILEGES ON DATABASE llama3_db TO postgres;" postgres

echo "Database initialized successfully!" 