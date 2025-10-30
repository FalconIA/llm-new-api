FROM oven/bun:latest AS builder

ARG BUN_CONFIG_REGISTRY
ENV BUN_CONFIG_REGISTRY=${BUN_CONFIG_REGISTRY}

WORKDIR /build
COPY web/package.json .
COPY web/bun.lock .
RUN bun install
COPY ./web .
COPY ./VERSION .
RUN DISABLE_ESLINT_PLUGIN='true' VITE_REACT_APP_VERSION=$(cat VERSION) bun run build

FROM golang:alpine AS builder2
ENV GO111MODULE=on CGO_ENABLED=0

ARG TARGETOS
ARG TARGETARCH
ARG GOPROXY
ENV GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} GOPROXY=${GOPROXY}

WORKDIR /build

ADD go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=builder /build/dist ./web/dist
RUN go build -ldflags "-s -w -X 'github.com/QuantumNous/new-api/common.Version=$(cat VERSION)'" -o new-api

FROM alpine

RUN apk upgrade --no-cache \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates

ENV TZ=Asia/Shanghai \
    TIKTOKEN_CACHE_DIR=/tiktoken-cache

ADD docker/one-api_tiktoken-cache.tar.gz /

COPY --from=builder2 /build/new-api /

EXPOSE 3000
WORKDIR /data
VOLUME /data
ENTRYPOINT ["/new-api"]
