# syntax=docker/dockerfile:1.4
# fix.Dockerfile for kubernetes-136226 (T3 stub with goproxy deps, FIX state)
FROM golang:1.24
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1
WORKDIR /work
COPY go.mod ./
COPY effective_alloc_fixed.go ./effective_alloc.go
COPY verified_test.go ./
RUN go mod tidy 2>&1 | tail -10
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=60s ."]
