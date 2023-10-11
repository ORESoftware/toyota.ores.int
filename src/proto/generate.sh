#!/usr/bin/env bash

dir_name="$(cd "$(dirname "${BASH_SOURCE}")" && PWD)"

protoc --go_out="$dir_name" --go-grpc_out="$dir_name" --proto_path="$dir_name" "${dir_name}/myservice.proto"
