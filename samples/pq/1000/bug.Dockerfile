FROM gonb-pq-1000-bug:latest
WORKDIR /work/upstream
RUN rm -f *_test.go
COPY verified_test.go ./pq_1000_race_test.go
COPY verified_test_helper.go ./pq_1000_helper.go
RUN sed -i 's,//go:build !pqfix,,' pq_1000_helper.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
