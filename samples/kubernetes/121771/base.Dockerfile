# base.Dockerfile for kubernetes-121771 (modules mode)
FROM golang:1.21
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates patch && rm -rf /var/lib/apt/lists/*
ENV GOPATH=/go GO111MODULE=on GOPROXY=https://goproxy.cn,direct GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1
RUN mkdir -p /work/upstream
COPY source_full/ /work/upstream/
COPY fix.diff /tmp/fix.diff
WORKDIR /work/upstream
# Revert fix.diff to get BUG state
RUN if [ -f /tmp/fix.diff ] && grep -q "^diff --git" /tmp/fix.diff; then \
      rm -rf .git 2>/dev/null; \
      git init --quiet && git add -A && \
      git -c user.email=x@x -c user.name=x commit -m base -q && \
      (git apply -R --whitespace=nowarn /tmp/fix.diff 2>&1 || patch -p1 -R < /tmp/fix.diff 2>&1 || echo "WARNING: fix.diff revert issues"); \
    fi; \
    echo "go.mod present:"; test -f go.mod && echo "YES" || echo "NO_MISSING"
RUN go mod download 2>&1 | tail -5 || true
