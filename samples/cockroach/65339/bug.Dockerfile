# bug.Dockerfile for cockroach-65339 (in-place BUG state)
FROM gonb-cockroach-65339-base:latest

WORKDIR /work/upstream/pkg/util/tracing

# Remove upstream .go files (keep only our stubs + race test)
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'crdbspan.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (unfixed / bug versions)
COPY crdbspan.go ./crdbspan.go
COPY verified_test.go ./verified_test.go

# Verify compilation
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'Test65339Race' .
