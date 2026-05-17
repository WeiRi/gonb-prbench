# syntax=docker/dockerfile:1.4
FROM golang:1.24
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off CGO_ENABLED=1
RUN apt-get update && apt-get install -y --no-install-recommends git patch ca-certificates && rm -rf /var/lib/apt/lists/*
RUN mkdir -p /root/.ssh && ssh-keyscan -t rsa,ed25519 github.com >> /root/.ssh/known_hosts 2>/dev/null
RUN --mount=type=ssh git clone git@github.com:OliveTin/OliveTin.git /work/upstream
WORKDIR /work/upstream
RUN git checkout --detach d3cd876eec1b8e0a4cfb8e857ede452230b8b2d5
WORKDIR /work/upstream/service
RUN --mount=type=ssh go mod download
