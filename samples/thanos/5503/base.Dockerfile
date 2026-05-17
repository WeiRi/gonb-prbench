# base.Dockerfile for thanos-5503 (Recipe A in-place)
FROM golang:1.25
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates patch && rm -rf /var/lib/apt/lists/*
ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB=off CGO_ENABLED=1 GOTOOLCHAIN=auto
COPY src /work/upstream
WORKDIR /work/upstream
# Handle pre-modules repos: create go.mod if missing (best-effort)
RUN (if [ ! -f go.mod ]; then go mod init github.com/thanos-io/thanos 2>/dev/null || true; fi; \
     go mod tidy 2>&1 | tail -3 || true; \
     go mod download 2>&1 | tail -5 || true) || true
