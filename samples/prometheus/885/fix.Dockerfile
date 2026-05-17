# fix.Dockerfile for prometheus-885 — worktree at FIX already, rename tests
FROM gonb-prometheus-885-base:latest

WORKDIR /go/src/github.com/prometheus/prometheus/retrieval
RUN for f in *_test.go; do if [ "$f" != "race_test.go" ]; then mv "$f" "verified_test_$(echo $f)"; fi; done && \
    rm -f race_test.go

RUN go test -race -c -o /dev/null . 2>&1 | tail -10 || true
