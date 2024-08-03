#!/usr/bin/env bash
set -x
set -eo pipefail

if ! [ -x "$(command -v mysql)" ]; then
    echo >&2 "Error: mysql is not installed."
    exit 1
fi

if ! [ -x "$(command -v mysqladmin)" ]; then
    echo >&2 "Error: mysqladmin is not installed."
    exit 1
fi

# Check if a custom root password has been set, otherwise default to 'password'
MYSQL_ROOT_PASSWORD="${MYSQL_ROOT_PASSWORD:=password}"
# Check if a custom user has been set, otherwise default to 'user'
MYSQL_USER="${MYSQL_USER:=user}"
# Check if a custom password has been set, otherwise default to 'password'
MYSQL_PASSWORD="${MYSQL_PASSWORD:=password}"
# Check if a custom database name has been set, otherwise default to 'BotDB'
MYSQL_DATABASE="${MYSQL_DATABASE:=BotDB}"
# Check if a custom port has been set, otherwise set to '3306'
MYSQL_PORT="${MYSQL_PORT:=3306}"
# Check if a custom host has been set, otherwise default to 'localhost'
DB_HOST="${DB_HOST:=localhost}"

# Debugging: Print environment variables
echo "MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}"
echo "MYSQL_USER=${MYSQL_USER}"
echo "MYSQL_PASSWORD=${MYSQL_PASSWORD}"
echo "MYSQL_DATABASE=${MYSQL_DATABASE}"
echo "MYSQL_PORT=${MYSQL_PORT}"
echo "DB_HOST=${DB_HOST}"

# Launch MySQL using docker
sudo docker run --rm \
 -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
 -e MYSQL_USER=${MYSQL_USER} \
 -e MYSQL_PASSWORD=${MYSQL_PASSWORD} \
 -e MYSQL_DATABASE=${MYSQL_DATABASE} \
 -p "${MYSQL_PORT}":3306  \
 -d mysql:latest --skip-networking=off

# Wait for MySQL to be ready
echo "Waiting for MySQL to initialize..."

sleep 5
until mysqladmin ping -h "${DB_HOST}" -P "${MYSQL_PORT}" --silent; do
    >&2 echo "MySQL is still unavailable - sleeping"
    sleep 1
done

>&2 echo "MySQL is up and running on port ${MYSQL_PORT}!"

# Create the database and run migrations using MySQL CLI
mysql -h "${DB_HOST}" -P "${MYSQL_PORT}" -u root -e "CREATE DATABASE IF NOT EXISTS ${MYSQL_DATABASE};"

# Assuming you have a directory with your migration SQL files
MIGRATIONS_DIR="./migrations"

# Check if migrations directory exists and is not empty
if [ ! -d "${MIGRATIONS_DIR}" ]; then
    echo >&2 "Error: Migrations directory '${MIGRATIONS_DIR}' does not exist."
    exit 1
fi

if [ -z "$(ls -A ${MIGRATIONS_DIR})" ]; then
    echo >&2 "Error: Migrations directory '${MIGRATIONS_DIR}' is empty."
    exit 1
fi

for migration in "${MIGRATIONS_DIR}"/*.sql; do
    >&2 echo "Running migration file ${migration}"
    mysql -h "${DB_HOST}" -P "${MYSQL_PORT}" -u root "${MYSQL_DATABASE}" < "${migration}"
done

# Seeding values to our DB
SEED_DIR="./seeding"
for seed in "${SEED_DIR}"/*.sql; do
    >&2 echo "Running seed file ${seed}"
    mysql -h "${DB_HOST}" -P "${MYSQL_PORT}" -u root "${MYSQL_DATABASE}" < "${seed}"
done
