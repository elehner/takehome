# syntax=docker/dockerfile:1

## Build the image
FROM golang:1.19.2-bullseye AS build

ENV APP_HOME /src/takehomeserver
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN mkdir "users"
COPY users/*.go ./users/
RUN mkdir "images"
COPY images/*.go ./images/
RUN go build -o /takehome-server

## Deploy the server
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /takehome-server /takehome-server

EXPOSE 8080
USER nonroot:nonroot
CMD [ "/takehome-server" ]