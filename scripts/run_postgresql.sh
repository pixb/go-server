#!/bin/bash

set -e

# Change to repo root
cd "$(dirname "$0")/../"

./build/go-server --driver postgresql --dsn "host=localhost port=5432 user=root password=123456 dbname=goserver sslmode=disable"
