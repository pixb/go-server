#!/bin/bash

./build/go-server --driver mysql --dsn "root:123456@tcp(localhost:3306)/goserver?charset=utf8mb4&parseTime=True&loc=Local"
