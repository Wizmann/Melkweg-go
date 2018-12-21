#!/bin/bash

echo "init golang/protobuf"

GIT_TAG="v1.2.0" # change as needed
go get -d -u github.com/golang/protobuf/protoc-gen-go
git -C "$(go env GOPATH)"/src/github.com/golang/protobuf checkout $GIT_TAG
go install github.com/golang/protobuf/protoc-gen-go

echo "init logurs"
go get github.com/sirupsen/logrus
go get -u golang.org/x/crypto/ssh/terminal

echo "done"
