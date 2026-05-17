FROM gonb-grpc-go-2074-base-v2:latest
ENV GO111MODULE=off
WORKDIR /go/src/google.golang.org/grpc/transport
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./2074_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
