FROM gonb-go-ethereum-20421-bug:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /work/upstream/p2p/enode
RUN rm -f *_test.go 2>/dev/null || true
COPY verified_test.go ./ge_20421_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
