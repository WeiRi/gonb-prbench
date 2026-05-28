# syntax=docker/dockerfile:1.4
# bug.Dockerfile for kubernetes-94537 — full upstream clone at bug commit, drops test into upstream pkg
FROM golang:1.17
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates patch && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /root/.ssh && ssh-keyscan -t rsa,ed25519 github.com >> /root/.ssh/known_hosts 2>/dev/null
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1

# === Full upstream at bug commit ===
RUN --mount=type=ssh git clone --depth=200 git@github.com:kubernetes/kubernetes.git /work/upstream
WORKDIR /work/upstream
RUN --mount=type=ssh git fetch --depth=200 origin a84419f02766c2b8d6ce4a887f6fe45c6751d46e && git checkout --detach a84419f02766c2b8d6ce4a887f6fe45c6751d46e
RUN --mount=type=ssh go mod download 2>&1 | tail -10 || true

COPY verified_test.go /work/upstream/staging/src/k8s.io/legacy-cloud-providers/azure/cache/94537_handcrafted_race_test.go

WORKDIR /work/upstream/staging/src/k8s.io/legacy-cloud-providers/azure/cache
# NO CMD — race trigger command is in README
