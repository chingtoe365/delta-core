FROM golang:1.22-alpine

# pre-download packages so only if mod file changes will it skip using cache
RUN mkdir /app-build
COPY ./src/go.mod /app-build/ 
COPY ./src/go.sum /app-build/ 
WORKDIR /app-build
RUN go mod download

RUN mkdir /app
COPY ./src /app
WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

COPY ./docker/scripts/start.sh /start.sh
RUN chmod +x /start.sh

ARG ENV_FILE=${ENV_FILE}
COPY $ENV_FILE /envs/.env