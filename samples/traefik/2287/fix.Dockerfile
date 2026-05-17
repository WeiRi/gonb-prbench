FROM gonb-traefik-2287-base:latest

WORKDIR /go/src/github.com/containous/traefik
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /healthcheck\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    git init --quiet && git add -A 2>/dev/null && git -c user.email=x@x -c user.name=x commit -m b -q 2>/dev/null && \
    git apply --whitespace=nowarn /tmp/fix_prod.diff

WORKDIR /go/src/github.com/containous/traefik/healthcheck
RUN rm -f race_2287_test.go && \
    for f in *_test.go; do mv "$f" "verified_test_$(echo $f)"; done 2>/dev/null

COPY verified_test.go ./traefik_2287_race_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
