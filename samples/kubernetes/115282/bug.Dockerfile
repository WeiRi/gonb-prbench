# syntax=docker/dockerfile:1.4
FROM golang:1.22
ENV GOPROXY=off GOSUMDB=off GOFLAGS=-mod=mod CGO_ENABLED=1
WORKDIR /work
COPY go.mod ./
COPY config.go ./
COPY verified_test.go ./
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true
CMD ["sh","-c","go test -race -vet=off -count=5 -timeout=60s ."]
