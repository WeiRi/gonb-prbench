FROM gonb-kubernetes-132063-base-v2:latest
WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /counter.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    (git apply --whitespace=nowarn /tmp/fix_prod.diff || (git init --quiet && git add -A && git -c user.email=x@x -c user.name=x commit -m b -q && git apply --whitespace=nowarn /tmp/fix_prod.diff))
WORKDIR /work/upstream/staging/src/k8s.io/component-base/metrics
RUN find . -maxdepth 1 -name "*_test.go" -exec sh -c 'mv "$1" "verified_test_$(basename $1)"' _ {} \; 2>/dev/null || true
COPY verified_test.go ./kubernetes_132063_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
