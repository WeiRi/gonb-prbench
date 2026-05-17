FROM gonb-thanos-5972-bug:latest
WORKDIR /work/upstream/pkg/query
RUN rm -f *_test.go 2>/dev/null || true
COPY verified_test.go ./thanos_5972_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
