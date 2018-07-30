#!/bin/bash

# github.com/google/protobuf
# github.com/golang/protobuf
# github.com/golang/protobuf/protoc-gen-go
# github.com/golang/grpc

# run at ./stormtf/pb
protoc -I storm feature.proto --go_out=plugins=grpc:./stormtf --proto_path="$HOME/go/src"