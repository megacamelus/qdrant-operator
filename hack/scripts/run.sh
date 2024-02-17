#!/bin/sh -x

go run -ldflags="${GOLDFLAGS}" cmd/main.go run --leader-election=false --zap-devel