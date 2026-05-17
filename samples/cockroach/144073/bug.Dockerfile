# bug.Dockerfile for cockroach-144073 (in-place BUG state)
FROM gonb-cockroach-144073-base:latest

WORKDIR /work/upstream/pkg/sql/catalog/lease

# Remove upstream .go files (keep only our stubs + race test)
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'descriptor_state.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (unfixed / bug versions)
COPY descriptor_state.go ./descriptor_state.go
COPY verified_test.go ./verified_test.go

# Verify compilation
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'Test144073Race' .
