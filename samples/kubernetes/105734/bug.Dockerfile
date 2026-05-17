# bug.Dockerfile for kubernetes-105734 (in-place BUG state)
FROM gonb-kubernetes-105734-base:latest

WORKDIR /work/upstream/staging/src/k8s.io/apiserver/pkg/server/httplog

# Remove upstream .go files (keep only our stubs + race test)
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'httplog.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (unfixed / bug versions)
COPY httplog.go ./httplog.go
COPY verified_test.go ./verified_test.go

# Verify compilation
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'TestRace_105734_HttplogConcurrent' .
