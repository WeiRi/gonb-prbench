# bug.Dockerfile for etcd-5505 (in-place BUG state)
FROM gonb-etcd-5505-base:latest

WORKDIR /go/src/github.com/coreos/etcd/etcdserver/api/v3rpc

# Remove upstream .go files (keep only our stubs + race test)
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'watch.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (unfixed / bug versions)
COPY watch.go ./watch.go
COPY verified_test.go ./verified_test.go

# Verify compilation
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'TestRace_PR5505_CtrlChanClose' .
