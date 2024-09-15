#!/bin/sh

cd $(dirname "$0")
go run ../cmd/central-auth-service/main.go
