#!/usr/bin/env bash


echo "build dir2http.exe"
go build dir2http.go

# echo "build dir2http"
# $env:GOOS="linux" 
# go build dir2http.go
# $env:GOOS=""


rm -ErrorAction SilentlyContinue wasmclockcal.wasm

$env:GOOS = "js" 
$env:GOARCH = "wasm" 
echo "go build -o clientdata/wasmclockcal.wasm wasmclockcal.go"
go build -o wasmclockcal.wasm wasmclockcal.go
$env:GOOS = ""
$env:GOARCH = ""


