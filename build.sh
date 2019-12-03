#!/usr/bin/env bash

go build dir2http.go

rm wasmclockcal.wasm

echo "GOOS=js GOARCH=wasm go build -o wasmclockcal.wasm"
GOOS=js GOARCH=wasm go build -o wasmclockcal.wasm
