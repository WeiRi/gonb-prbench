FROM gonb-nats-server-540-base:latest
ENV GO111MODULE=off
RUN rm -rf /work 2>/dev/null || true
RUN ln -sf /go/src/github.com/nats-io/nats-server /go/src/github.com/nats-io/gnatsd
WORKDIR /go/src/github.com/nats-io/nats-server
COPY fix_prod.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff
WORKDIR /go/src/github.com/nats-io/nats-server/server
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./540_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
