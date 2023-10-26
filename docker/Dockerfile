FROM golang:1.19-alpine

RUN mkdir /app

ADD ./src /app

WORKDIR /app

COPY ./docker/scripts/start.sh /start.sh
RUN chmod +x /start.sh

ARG ENV_FILE=${ENV_FILE}
COPY $ENV_FILE /envs/.env

RUN go mod download

#CMD ["/app/main"]