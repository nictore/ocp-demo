#!/bin/bash

GOOGLE_APIS=${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis

for proto in user product order account
do
  echo "Removing ${proto}/*.go"
  rm -f ${proto}/*.go
  echo "Generating ${proto}/${proto}.go"
  protoc -I . -I ${GOOGLE_APIS} --go_out=paths=source_relative,plugins=grpc:. ${proto}/${proto}.proto
  echo "Generating ${proto}/${proto}.gw.go"
  protoc -I . -I ${GOOGLE_APIS} --grpc-gateway_out=paths=source_relative:. ${proto}/${proto}.proto
done
