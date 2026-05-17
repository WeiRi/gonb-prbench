# bug.Dockerfile for thanos-5503 (Recipe A in-place)
FROM gonb-thanos-5503-base:latest

WORKDIR /work/upstream/pkg/compact
# Keep existing tests (admission gate handles cleanup at runtime)

COPY verified_test_inplace.go ./thanos_5503_race_test.go

# Compile check
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true

CMD go test -race -vet=off -count=10 -timeout=300s .
