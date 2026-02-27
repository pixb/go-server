#!/bin/bash

set -e

# Change to repo root
cd "$(dirname "$0")/../"

./build/go-server --driver sqlite --data ./data/
