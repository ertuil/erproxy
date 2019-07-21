#!/bin/sh
cd ..
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o docker/erproxy erproxy
cd docker
docker build -t ertuil/erproxy:latest .