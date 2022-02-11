#!/bin/sh

GOOS=linux GOARCH=amd64 go build -ldflags "$flags" -o main main.go