# fix.Dockerfile for nomad-14121 — apply fix + preserve upstream tests
FROM gonb-nomad-14121-base:latest

WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    cd /work/upstream && (git init --quiet 2>/dev/null; git add -A 2>/dev/null; git -c user.email=x@x -c user.name=x commit -m b -q 2>/dev/null; true) && \
    git apply --whitespace=nowarn /tmp/fix_prod.diff

WORKDIR /work/upstream/nomad/deploymentwatcher
RUN mv deployments_watcher_test.go verified_test_dw_test.go && \
    mv testutil_test.go verified_test_dwutil_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
