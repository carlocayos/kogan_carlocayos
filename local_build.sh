#!/usr/bin/env bash

# Build the binary and save to bin folder
go mod download
go build -o ./bin/myapp app/*