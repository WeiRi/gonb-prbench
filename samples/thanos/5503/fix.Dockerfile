# fix.Dockerfile for thanos-5503 (Recipe A in-place)
FROM gonb-thanos-5503-base:latest

# Apply fix.diff (production changes only) on top of workspace (best-effort)
COPY fix.diff /tmp/fix.diff
WORKDIR /work/upstream
# Filter out test/md/bazel files, then apply fix
RUN (awk 'BEGIN{p=0} /^diff --git/{if ($$0 !~ /_test\.go/ && $$0 !~ /\.md$$/ && $$0 !~ /BUILD\.bazel/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff; \
     git init -q 2>/dev/null; git add -A 2>/dev/null; \
     git -c user.email=x@x -c user.name=x commit -m base -q 2>/dev/null; \
     git apply --whitespace=nowarn /tmp/fix_prod.diff 2>/dev/null || \
     patch -p1 -f < /tmp/fix.diff 2>/dev/null || \
     echo 'WARNING: fix.diff apply had issues') || true

WORKDIR /work/upstream/pkg/compact
# Keep existing tests (admission gate handles cleanup at runtime)

COPY verified_test_inplace.go ./thanos_5503_race_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true

CMD go test -race -vet=off -count=10 -timeout=300s .
