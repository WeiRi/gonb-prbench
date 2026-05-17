# bug.Dockerfile for nomad-14119 (Recipe A in-place)
FROM gonb-nomad-14119-base:latest

WORKDIR /work/upstream/client
# Keep existing tests (admission gate handles cleanup at runtime)

COPY verified_test_inplace.go ./nomad_14119_race_test.go

# Compile check
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true

CMD go test -race -vet=off -count=10 -timeout=300s .
