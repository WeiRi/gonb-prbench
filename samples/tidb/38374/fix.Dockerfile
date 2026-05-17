# fix.Dockerfile for tidb-38374 — apply fix.diff (BUG state→FIX) + use FIX-API test
FROM gonb-tidb-38374-base:latest

WORKDIR /work/upstream
RUN rm -rf .git
COPY fix.diff /tmp/fix.diff
# Apply only the .go portions (skip BUILD.bazel and _test.go)
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 !~ /_test\.go/ && $0 !~ /BUILD\.bazel/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    git init --quiet && git add -A 2>/dev/null && git -c user.email=x@x -c user.name=x commit -m b -q 2>/dev/null && \
    git apply --whitespace=nowarn /tmp/fix_prod.diff

WORKDIR /work/upstream/sessionctx/variable
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null

COPY verified_test_fix.go ./tidb_38374_race_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
