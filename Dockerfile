# syntax=docker/dockerfile:1.14-labs@sha256:bbd3af5735ce9921cfa2a8b59d1006445333e161cb4582c3d88aad0730d9b207
FROM cgr.dev/chainguard/wolfi-base:latest@sha256:7b05dd6ccd1abdc790b54c64e8ab0bc0d994c39ea04bbc115a021562d5b42bb7 as base
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