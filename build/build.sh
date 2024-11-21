#!/bin/bash

cd "src"

echo "Building contor system..."

go build -o ../ contorsystem.go

if [ $? -eq 0 ]; then
    echo "Build successful. Executable created as 'main' in the root directory."
else
    echo "Build failed."
    exit 1
fi