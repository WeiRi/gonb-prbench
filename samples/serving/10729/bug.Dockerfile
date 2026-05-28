# syntax=docker/dockerfile:1.4
# bug.Dockerfile for serving-10729 — full upstream clone at bug commit, drops test into upstream pkg
FROM golang:1.17
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates patch && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /root/.ssh && ssh-keyscan -t rsa,ed25519 github.com >> /root/.ssh/known_hosts 2>/dev/null
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1

# === Full upstream at bug commit ===
RUN --mount=type=ssh git clone --depth=200 git@github.com:knative/serving.git /work/upstream
WORKDIR /work/upstream
RUN --mount=type=ssh git fetch --depth=200 origin c9fc401a8a484fd2a667e5f2e008b5650d62b989 && git checkout --detach c9fc401a8a484fd2a667e5f2e008b5650d62b989
RUN --mount=type=ssh go mod download 2>&1 | tail -10 || true

COPY verified_test.go /work/upstream/pkg/activator/handler/10729_handcrafted_race_test.go

WORKDIR /work/upstream/pkg/activator/handler
# NO CMD — race trigger command is in README
