#!/bin/sh
# start.sh

set -e

cmd="$@"

# wait for clickhouse to be ready
# until ./clickhouse client -h $DB_HOST -u $DB_USER -d $DB_NAME; do
#   >&2 echo "Clickhouse is unavailable - sleeping"
#   sleep 1
# done

# >&2 echo "Clickhouse is up - executing command"

# make migrations
#goose -dir ./migrations up

# run tests
go test -v ./...

exec $cmd
