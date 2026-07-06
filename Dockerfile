# syntax=docker/dockerfile:1.25-labs@sha256:4426b5e269e36911b94fb79cf67f1fd7155ef11b2bbc8ab23cbfcbc97130efe9
FROM cgr.dev/chainguard/wolfi-base:latest@sha256:42df77a9974d6ec8b17a5ee8bc23b532600a44d705acef2409e0933c1251b45f as base
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