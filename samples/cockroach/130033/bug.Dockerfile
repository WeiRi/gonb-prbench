# bug.Dockerfile for cockroach-130033 (in-place BUG state)
FROM gonb-cockroach-130033-base:latest

WORKDIR /work/upstream/pkg/server/diagnostics

# Remove upstream .go files (keep only our stubs + race test)
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'reporter.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (unfixed / bug versions)
COPY reporter.go ./reporter.go
COPY verified_test.go ./verified_test.go

# Verify compilation
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'Test130031Race' .
