#!/bin/bash

# Build the Go application
go build cmd/main.go

# Check if the build was successful
if [ $? -eq 0 ]; then
    # Run the generated executable
    ./main
else
    echo "Build failed. Unable to run the application."
fi