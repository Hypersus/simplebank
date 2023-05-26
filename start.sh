#!/bin/sh
# This script is used to create the database and the tables for the docker container
# Include this shell script in the Dockerfile and RUN it before the ENTRYPOINT

set -e

echo "migrating database"
source /app/app.env
export DB_SOURCE="postgres://postgres:aA5\!bB4\#bK0\!pU4\#@192.168.0.212:5432/simple_bank"

echo "DB_SOURCE: ${DB_SOURCE}"
/app/migrate -path /app/migration -database ${DB_SOURCE} --verbose up
echo "migrated database"

echo "starting server"
exec "$@"
# while true; do sleep 30; done;