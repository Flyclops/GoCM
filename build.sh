#!/bin/bash
export GOPATH=`pwd`

rm -rf bin pkg src
mkdir bin

go get github.com/alexjlockwood/gcm
go build -o bin/gocm main.go 
