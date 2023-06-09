# syntax=docker/dockerfile:1.3-labs

FROM debian:bullseye-slim as base
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
RUN useradd -r -u 999 -d /home/go-project-template go-project-template

FROM ghcr.io/acorn-io/images-mirror/golang:1.20 AS build
COPY / /src
WORKDIR /src
RUN \
  --mount=type=cache,target=/go/pkg \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o bin/go-project-template main.go

FROM base AS goreleaser
COPY go-project-template /usr/local/bin/go-project-template
USER go-project-template

FROM base
COPY --from=build /src/bin/go-project-template /usr/local/bin/go-project-template
USER go-project-template