FROM gonb-kubernetes-129170-base-v5:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /custom\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    patch -p1 -i /tmp/fix_prod.diff
WORKDIR /work/upstream/staging/src/k8s.io/apimachinery/pkg/runtime/serializer/cbor/internal/modes
RUN find . -maxdepth 1 -name "*_test.go" -exec sh -c 'mv "$1" "verified_test_$(basename $1)"' _ {} \; 2>/dev/null || true
COPY verified_test.go ./129170_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
