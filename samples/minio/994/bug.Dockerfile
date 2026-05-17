FROM gonb-minio-994-base-v2:latest
WORKDIR /work/upstream/pkg/contentdb
RUN for f in *_test.go; do mv "$f" "verified_test_$(echo $f)"; done 2>/dev/null
COPY verified_test.go ./minio_994_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
