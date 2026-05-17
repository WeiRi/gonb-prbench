# fix.Dockerfile for etcd-5739 (in-place FIX state)
FROM gonb-etcd-5739-base:latest

WORKDIR /go/src/github.com/coreos/etcd
COPY fix.diff /tmp/fix.diff
# Apply fix_prod.diff (filter out test/md/bazel files)
RUN awk 'BEGIN{p=0} /^diff --git/{if ($$0 !~ /_test\.go/ && $$0 !~ /\.md$$/ && $$0 !~ /BUILD\.bazel/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    git apply --whitespace=nowarn /tmp/fix_prod.diff 2>/dev/null || patch -p1 < /tmp/fix.diff 2>/dev/null || echo "WARNING: fix apply issues"

WORKDIR /go/src/github.com/coreos/etcd/etcdserver

# Remove upstream .go files
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'apply_auth.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (fixed variants renamed to original names)
COPY apply_auth_fixed.go ./apply_auth.go
COPY verified_test.go ./verified_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'TestRace_PR5739_AuthApplierUserUnlocked' .
