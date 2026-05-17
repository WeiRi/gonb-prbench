FROM gonb-teleport-65910-bug:latest
WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /handler\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    git apply --whitespace=nowarn /tmp/fix_prod.diff
WORKDIR /work/upstream/lib/srv/app/gcp
RUN for f in *_test.go; do mv "$f" "verified_test_$(echo $f)"; done 2>/dev/null
COPY verified_test.go ./teleport_65910_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
