FROM gonb-moby-29893-base-v6:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /go/src/github.com/docker/docker/pkg/plugins
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./29893_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
