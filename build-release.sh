#!/bin/bash

# Set the custom binary name
BINARY_NAME="pixelita"

# Build the Go source code with optimizations
go build -o "$BINARY_NAME" -ldflags "-s -w" -gcflags "all=-N -l" ./cmd/main.go

# Check if the build was successful
if [ $? -eq 0 ]; then

    # Check if the compression was successful
    if [ $? -eq 0 ]; then
        echo "Build successful! The binary '$BINARY_NAME' is ready."
    else
        echo "Compression failed."
    fi
else
    echo "Build failed."
fi
