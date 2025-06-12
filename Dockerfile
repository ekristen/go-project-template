# syntax=docker/dockerfile:1.16-labs@sha256:bb5e2b225985193779991f3256d1901a0b3e6a0b284c7bffa0972064f4a6d458
FROM cgr.dev/chainguard/wolfi-base:latest@sha256:08a4c4fc8583c217c853fda751f08495530d105c361b714f6d33ae3edb5ec11c as base
ARG PROJECT_NAME=go-project-template
RUN apk add --no-cache ca-certificates
RUN addgroup -S ${PROJECT_NAME} && adduser -S ${PROJECT_NAME} -G ${PROJECT_NAME}

FROM ghcr.io/acorn-io/images-mirror/golang:1.21@sha256:856073656d1a517517792e6cdd2f7a5ef080d3ca2dff33e518c8412f140fdd2d AS build
ARG PROJECT_NAME=go-project-template
COPY / /src
WORKDIR /src
RUN \
  --mount=type=cache,target=/go/pkg \
  --mount=type=cache,target=/root/.cache/go-build \
  go build -o bin/${PROJECT_NAME} main.go

FROM base AS goreleaser
ARG PROJECT_NAME=go-project-template
COPY ${PROJECT_NAME} /usr/local/bin/${PROJECT_NAME}
USER ${PROJECT_NAME}

FROM base
ARG PROJECT_NAME=go-project-template
COPY --from=build /src/bin/${PROJECT_NAME} /usr/local/bin/${PROJECT_NAME}
USER ${PROJECT_NAME}