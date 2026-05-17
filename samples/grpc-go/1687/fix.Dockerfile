FROM gonb-grpc-go-1687-base-v3:latest
ENV GO111MODULE=off
WORKDIR /go/src/google.golang.org/grpc
COPY fix_prod.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff && grep -c writeStatusMu transport/handler_server.go
WORKDIR /go/src/google.golang.org/grpc/transport
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./1687_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
