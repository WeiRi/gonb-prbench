# bug.Dockerfile for kubernetes-107452 (in-place BUG state)
FROM gonb-kubernetes-107452-base:latest

WORKDIR /work/upstream/staging/src/k8s.io/apiserver/pkg/server/filters

# Remove upstream .go files (keep only our stubs + race test)
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'timeout.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (unfixed / bug versions)
COPY timeout.go ./timeout.go
COPY verified_test.go ./verified_test.go

# Verify compilation
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'TestRace_107452_HeaderMutation' .
