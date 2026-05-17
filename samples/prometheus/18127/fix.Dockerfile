FROM gonb-prometheus-18127-bug:latest
ENV GOWORK=off
WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /scrape\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    (git apply --whitespace=nowarn /tmp/fix_prod.diff || (git init --quiet && git add -A && git -c user.email=x@x -c user.name=x commit -m b -q && git apply --whitespace=nowarn /tmp/fix_prod.diff))
WORKDIR /work/upstream/scrape
RUN for f in *_test.go; do mv "$f" "verified_test_$(echo $f)"; done 2>/dev/null
COPY verified_test.go ./prometheus_18127_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
