FROM gonb-go-ethereum-20421-bug:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /iter\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    patch -p1 -i /tmp/fix_prod.diff
WORKDIR /work/upstream/p2p/enode
RUN rm -f *_test.go 2>/dev/null || true
COPY verified_test.go ./ge_20421_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
