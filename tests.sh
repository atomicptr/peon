#!/usr/bin/env bash
set -e

go test ./... -json -coverpkg=./... -coverprofile=coverage.out | tparse
go tool cover -html=coverage.out -o coverage.html
go tool cover -func=coverage.out | awk '/total:/ {print "\n\tTotal Project Coverage: " $3}'

echo -ne "\n\t🥳 All tests completed successfully\n"
