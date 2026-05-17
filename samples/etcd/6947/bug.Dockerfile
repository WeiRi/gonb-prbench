# bug.Dockerfile for etcd-6947 (in-place BUG state)
FROM gonb-etcd-6947-base:latest

WORKDIR /go/src/github.com/coreos/etcd/proxy/grpcproxy/cache

# Remove upstream .go files (keep only our stubs + race test)
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'store.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (unfixed / bug versions)
COPY store.go ./store.go
COPY verified_test.go ./verified_test.go

# Verify compilation
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'TestRace_PR6947_SizeUnlocked' .
