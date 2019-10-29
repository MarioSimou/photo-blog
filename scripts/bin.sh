#!/bin/bash

scripts_folder=$(dirname "$(realpath $0)")
app="$scripts_folder/../cmd/authAPI/main.go"
build_path="$scripts_folder/../build/app"
env_path="$scripts_folder/../configs/.env"

go build -o $build_path $app