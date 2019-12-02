#!/usr/bin/env bash


rm wasmclockcal.wasm

echo "GOOS=js GOARCH=wasm go build -o wasmclockcal.wasm"
GOOS=js GOARCH=wasm go build -o wasmclockcal.wasm

go run dir2http.go
