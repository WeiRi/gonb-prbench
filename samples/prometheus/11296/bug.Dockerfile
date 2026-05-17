FROM gonb-prometheus-11296-bug:latest
WORKDIR /work/upstream/tsdb/chunkenc
RUN for f in *_test.go; do mv "$f" "verified_test_$(echo $f)"; done 2>/dev/null
COPY verified_test.go ./prometheus_11296_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
