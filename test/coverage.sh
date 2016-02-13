#!/bin/bash
cd "$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
go get github.com/ChrisTrenkamp/goxpath/...
if [ $? != 0 ]; then
	exit 1
fi
go test >/dev/null 2>&1
if [ $? != 0 ]; then
	go test
	exit 1
fi
gometalinter --deadline=20s ../...
go test -coverprofile=coverage.out -coverpkg=github.com/ChrisTrenkamp/goxpath/... >/dev/null 2>&1
go tool cover -html=coverage.out -o coverage.html >/dev/null 2>&1
firefox coverage.html
rm coverage.out coverage.html
