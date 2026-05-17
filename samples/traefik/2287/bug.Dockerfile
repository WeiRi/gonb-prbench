FROM gonb-traefik-2287-base:latest

WORKDIR /go/src/github.com/containous/traefik/healthcheck
# Rename upstream tests to survive CLEAN; remove our outdated synthetic
RUN rm -f race_2287_test.go && \
    for f in *_test.go; do mv "$f" "verified_test_$(echo $f)"; done 2>/dev/null

COPY verified_test.go ./traefik_2287_race_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
