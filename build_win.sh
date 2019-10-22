#!/bin/bash
source env.sh
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o output_win/client.exe Melkweg/Client
