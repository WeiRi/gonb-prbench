# bug.Dockerfile for prometheus-885 — worktree at FIX, revert target.go
FROM gonb-prometheus-885-base:latest

WORKDIR /go/src/github.com/prometheus/prometheus
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /target\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_target_only.diff && \
    git init --quiet 2>/dev/null && git add -A 2>/dev/null && git -c user.email=x@x -c user.name=x commit -m b -q 2>/dev/null; \
    git apply -R --whitespace=nowarn /tmp/fix_target_only.diff

WORKDIR /go/src/github.com/prometheus/prometheus/retrieval
RUN for f in *_test.go; do if [ "$f" != "race_test.go" ]; then mv "$f" "verified_test_$(echo $f)"; fi; done && \
    rm -f race_test.go && ls *_test.go | head

RUN go test -race -c -o /dev/null . 2>&1 | tail -10 || true
