#!/bin/bash
sorurce .env

export MIGRATION_DSN="host=pg port=5432 dbname=$PG_DATABASSE_NAME user=$PG_USER password=$PG_PASSWORD sslmode=disable"

sleep 2 && goose -dir "${MIGRATION_DIR}" postgres "${MIGRATION_DSN}" up -v