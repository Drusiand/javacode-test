FROM golang:bookworm

LABEL maintainer="Andrey Ivanov"

RUN apt update && apt install git && apt install bash

RUN mkdir /app
WORKDIR /app

COPY . .
COPY config.env .


RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o /build cmd/app/main.go

EXPOSE 8080

CMD [ "/build" ]