FROM gonb-teleport-65910-bug:latest
WORKDIR /work/upstream/lib/srv/app/gcp
# Preserve handler_test.go (provides test helpers); rename to verified_test_*
RUN for f in *_test.go; do mv "$f" "verified_test_$(echo $f)"; done 2>/dev/null
COPY verified_test.go ./teleport_65910_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
