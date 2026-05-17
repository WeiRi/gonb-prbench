FROM gonb-etcd-13203-bug:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN patch -p1 -i /tmp/fix.diff
WORKDIR /work/upstream/client/v3
RUN find . -maxdepth 1 -name "*_test.go" -exec sh -c 'mv "$1" "verified_test_$(basename $1)"' _ {} \; 2>/dev/null || true
COPY verified_test_fixed.go ./13203_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
