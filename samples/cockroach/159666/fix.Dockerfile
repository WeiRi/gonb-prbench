# fix.Dockerfile for cockroach-159666 (in-place FIX state)
FROM gonb-cockroach-159666-base:latest

WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
# Apply fix_prod.diff (filter out test/md/bazel files)
RUN awk 'BEGIN{p=0} /^diff --git/{if ($$0 !~ /_test\.go/ && $$0 !~ /\.md$$/ && $$0 !~ /BUILD\.bazel/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    git apply --whitespace=nowarn /tmp/fix_prod.diff 2>/dev/null || patch -p1 < /tmp/fix.diff 2>/dev/null || echo "WARNING: fix apply issues"

WORKDIR /work/upstream/pkg/server/apiinternal

# Remove upstream .go files
RUN find . -maxdepth 1 -name '*.go' -a ! -name 'api_internal.go' -a ! -name 'verified_test.go' -delete 2>/dev/null || true

# Copy stub mock files (fixed variants renamed to original names)
COPY api_internal_fixed.go ./api_internal.go
COPY verified_test.go ./verified_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true

CMD go test -race -vet=off -count=20 -timeout=300s -run 'Test159666Race' .
