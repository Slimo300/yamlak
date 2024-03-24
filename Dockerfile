# syntax=docker/dockerfile:1

FROM golang:1.22-bookworm AS build

WORKDIR /app

COPY go.mod ./

COPY  . ./
RUN CGO_ENABLED=0 go build -o yamlak

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build app/yamlak /yamlak

ENTRYPOINT ["/yamlak"]