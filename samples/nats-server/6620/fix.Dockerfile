FROM gonb-nats-server-6620-base-v3:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN patch -p1 -i /tmp/fix.diff
WORKDIR /work/upstream/server
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test_fixed.go ./6620_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
