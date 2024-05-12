#!/bin/sh
# start.sh

set -e

cmd="$@"

# make migrations
goose -dir ./migrations up

# run tests
go test -v ./...

exec $cmd
