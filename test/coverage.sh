#!/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
go test -coverprofile=coverage.out -coverpkg=github.com/ChrisTrenkamp/goxpath/... 
go tool cover -html=coverage.out 
firefox coverage.out 
rm coverage.out
