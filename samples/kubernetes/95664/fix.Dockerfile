# syntax=docker/dockerfile:1.4
# fix.Dockerfile for kubernetes-95664 — bug + fix.diff applied
FROM golang:1.21
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates patch && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /root/.ssh && ssh-keyscan -t rsa,ed25519 github.com >> /root/.ssh/known_hosts 2>/dev/null
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1

# === Full upstream at bug commit ===
RUN --mount=type=ssh git clone --depth=200 git@github.com:kubernetes/kubernetes.git /work/upstream
WORKDIR /work/upstream
RUN --mount=type=ssh git fetch --depth=200 origin 7791bbe978d4dbcbb5694565c4a1c856bca48223 && git checkout --detach 7791bbe978d4dbcbb5694565c4a1c856bca48223
COPY fix.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff || patch -p1 < /tmp/fix.diff
RUN --mount=type=ssh go mod download 2>&1 | tail -10 || true

# === Race-triggering artefact in isolated sub-package ===
WORKDIR /work/pr2t-test
COPY go.mod ./
COPY verified_test.go ./
COPY *.go ./

WORKDIR /work
# NO CMD
