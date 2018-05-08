#!/bin/bash
protoc -I . feature.proto --go_out=plugins=grpc:../stormtf --proto_path="$HOME/go/src"