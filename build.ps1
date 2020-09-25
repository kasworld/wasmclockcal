#!/usr/bin/env bash

go build dir2http.go

rm -ErrorAction SilentlyContinue wasmclockcal.wasm

$env:GOOS="js" 
$env:GOARCH="wasm" 
Write-Output "go build -o clientdata/wasmclockcal.wasm wasmclockcal.go"
go build -o wasmclockcal.wasm wasmclockcal.go
$env:GOOS=""
$env:GOARCH=""


