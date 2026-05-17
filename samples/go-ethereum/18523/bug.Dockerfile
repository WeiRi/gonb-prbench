FROM gonb-go-ethereum-18523-base:latest
RUN rm -rf /work 2>/dev/null || true
WORKDIR /go/src/github.com/ethereum/go-ethereum/swarm/pss
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./18523_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
