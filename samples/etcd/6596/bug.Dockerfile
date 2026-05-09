# syntax=docker/dockerfile:1.4
# bug.Dockerfile for etcd-6596 — full upstream clone at bug commit
FROM golang:1.16
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates patch && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /root/.ssh && ssh-keyscan -t rsa,ed25519 github.com >> /root/.ssh/known_hosts 2>/dev/null
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1

# === Full upstream at bug commit ===
RUN --mount=type=ssh git clone --depth=200 git@github.com:etcd-io/etcd.git /work/upstream
WORKDIR /work/upstream
RUN --mount=type=ssh git fetch --depth=200 origin b8079b7fc05fa3266d47ab4dce3d874563951fcb && git checkout --detach b8079b7fc05fa3266d47ab4dce3d874563951fcb
RUN --mount=type=ssh go mod download 2>&1 | tail -10 || true

# === Race-triggering artefact in isolated sub-package ===
WORKDIR /work/pr2t-test
COPY go.mod ./
COPY verified_test.go ./
COPY *.go ./

WORKDIR /work
# NO CMD — race trigger command is in README
