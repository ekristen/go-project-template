# syntax=docker/dockerfile:1.15-labs@sha256:94edd5b349df43675bd6f542e2b9a24e7177432dec45fe3066bfcf2ab14c4355
FROM cgr.dev/chainguard/wolfi-base:latest@sha256:cdab2937f46d3485a80b1e3b4e4e0dffa510df62c1b2b372c4c047caefd71a2a as base
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