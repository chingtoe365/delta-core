#!/usr/bin/env sh

if [ $APP_ENV = development ]; 
then
    echo ">> Development mode"
    # go run /app/main.go
    gin run main.go --immediate --laddr :3000 --port 3000 --appPort 3001
else
    echo ">> Production mode"
    go build -o main /app/main.go && \
        /app/main
fi