FROM gonb-prometheus-18127-bug:latest
ENV GOWORK=off
WORKDIR /work/upstream/scrape
RUN for f in *_test.go; do mv "$f" "verified_test_$(echo $f)"; done 2>/dev/null
COPY verified_test.go ./prometheus_18127_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
