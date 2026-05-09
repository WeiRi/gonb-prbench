# syntax=docker/dockerfile:1.4
# bug.Dockerfile for pq-1000 — full upstream clone at bug commit
FROM golang:1.22
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates patch && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /root/.ssh && ssh-keyscan -t rsa,ed25519 github.com >> /root/.ssh/known_hosts 2>/dev/null
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1

# === Full upstream at bug commit ===
RUN --mount=type=ssh git clone --depth=200 git@github.com:lib/pq.git /work/upstream
WORKDIR /work/upstream
RUN --mount=type=ssh git fetch --depth=200 origin be2b75c2254d56a8565b84cc48390157d56d5d3a && git checkout --detach be2b75c2254d56a8565b84cc48390157d56d5d3a
RUN --mount=type=ssh go mod download 2>&1 | tail -10 || true

# === Race-triggering artefact in isolated sub-package ===
WORKDIR /work/pr2t-test
COPY go.mod ./
COPY verified_test.go ./
COPY *.go ./

WORKDIR /work
# NO CMD — race trigger command is in README
