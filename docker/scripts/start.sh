#!/usr/bin/env sh

if [ $APP_ENV = development ]; 
then
    echo ">> Development mode"
    go run /app/main.go
else
    echo ">> Production mode"
    go build -o main /app/main.go && \
        /app/main
fi