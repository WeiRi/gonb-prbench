FROM golang:1.24
RUN apt-get update && apt-get install -y --no-install-recommends git ca-certificates patch && rm -rf /var/lib/apt/lists/*
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off GOFLAGS= CGO_ENABLED=1 GOWORK=off
COPY --from=gonb-kubernetes-133587-bug:latest /work/upstream /work/upstream
WORKDIR /work/upstream
