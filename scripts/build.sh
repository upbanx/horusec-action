#!/usr/bin/env bash

export GO111MODULE=auto
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/horusec-commenter-linux-amd64 ./cmd/main