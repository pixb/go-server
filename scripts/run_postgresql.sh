#!/bin/bash

./build/go-server --driver postgresql --dsn "host=localhost port=5432 user=root password=123456 dbname=goserver sslmode=disable"
