FROM gonb-gofiber-4062-base-v3:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /work/upstream
RUN find . -maxdepth 1 -name "*_test.go" -exec sh -c 'mv "$1" "verified_test_$(basename $1)"' _ {} \; 2>/dev/null || true
COPY verified_test.go ./4062_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
