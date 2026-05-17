FROM gonb-nats-server-140-base:latest
ENV GO111MODULE=off
RUN rm -rf /work 2>/dev/null || true
WORKDIR /go/src/github.com/nats-io/gnatsd/server
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./140_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
