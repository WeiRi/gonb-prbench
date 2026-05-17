FROM golang:1.16
ENV CGO_ENABLED=1
WORKDIR /app
COPY verified_test.go ./consul_499_race_test.go
RUN go mod init ase/consul-499 2>/dev/null || true
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
