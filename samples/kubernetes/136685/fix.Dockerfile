# syntax=docker/dockerfile:1.4
# fix.Dockerfile for kubernetes-136685 (self-contained T3 stub, FIX state)
FROM golang:1.22
ENV GOPROXY=off GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1
WORKDIR /work
COPY go.mod ./
COPY helpers_fixed.go ./helpers.go
COPY verified_test.go ./
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=60s ."]
