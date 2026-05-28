# syntax=docker/dockerfile:1.4
# fix.Dockerfile for kubernetes-98653 — bug + fix.diff applied
FROM golang:1.17
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates patch && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /root/.ssh && ssh-keyscan -t rsa,ed25519 github.com >> /root/.ssh/known_hosts 2>/dev/null
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1

# === Full upstream at bug commit ===
RUN --mount=type=ssh git clone --depth=200 git@github.com:kubernetes/kubernetes.git /work/upstream
WORKDIR /work/upstream
RUN --mount=type=ssh git fetch --depth=200 origin 2783f2f76ec57b9831b91e6c4b35d35cee4345e7 && git checkout --detach 2783f2f76ec57b9831b91e6c4b35d35cee4345e7
COPY fix.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff || patch -p1 < /tmp/fix.diff
RUN --mount=type=ssh go mod download 2>&1 | tail -10 || true

COPY verified_test.go /work/upstream/staging/src/k8s.io/apimachinery/pkg/watch/98653_handcrafted_race_test.go

WORKDIR /work/upstream/staging/src/k8s.io/apimachinery/pkg/watch
# NO CMD
