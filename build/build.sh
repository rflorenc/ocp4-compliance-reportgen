#!/bin/bash

output_binary=${1:-reportgen}
goos=${2:-linux}
goarch=${3:-amd64}

GOOS=${goos} GOARCH=${goarch} go build -o ${output_binary} ../cmd/main.go
[[ $? == 0 ]] && echo "Built ${output_binary}" || exit 1
