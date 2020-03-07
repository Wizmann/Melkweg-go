#!/bin/bash
source env.sh
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o output_arm/client.exe Melkweg/Client
