#!/bin/bash
CGO_ENABLED=1 && CGO_LDFLAGS="-static" && go build -mod=vendor -ldflags="-s -w" -o bpfgen ./main.go

exit 0
