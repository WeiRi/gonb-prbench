FROM golang:1.22
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates patch && rm -rf /var/lib/apt/lists/*
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1
COPY --from=gonb-minio-994-bug:latest /work/upstream /work/upstream
WORKDIR /work/upstream
RUN rm -rf .git 2>/dev/null; go mod init github.com/minio/minio 2>/dev/null; go mod tidy 2>&1 | tail -3 || true
