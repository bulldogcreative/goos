#!/bin/bash

version=v1.0.0

env GOOS=windows GOARCH=amd64 go build -o dist/goos.exe cmd/main.go
cd dist
tar -zcvf windows-$version.tar.gz goos.exe
rm goos.exe
cd ..

env GOOS=linux GOARCH=amd64 go build -o dist/goos cmd/main.go
cd dist
tar -zcvf linux-$version.tar.gz goos
rm goos
cd ..
