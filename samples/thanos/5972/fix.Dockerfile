FROM gonb-thanos-5972-bug:latest
WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /endpointset\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    (git apply --whitespace=nowarn /tmp/fix_prod.diff || (git init --quiet && git add -A && git -c user.email=x@x -c user.name=x commit -m b -q && git apply --whitespace=nowarn /tmp/fix_prod.diff))
WORKDIR /work/upstream/pkg/query
RUN rm -f *_test.go 2>/dev/null || true
COPY verified_test.go ./thanos_5972_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
