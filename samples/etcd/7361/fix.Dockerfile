FROM gonb-etcd-7361-base-v7:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /go/src/github.com/coreos/etcd
COPY fix.diff /tmp/fix.diff
RUN patch -p1 -i /tmp/fix.diff
WORKDIR /go/src/github.com/coreos/etcd/proxy/tcpproxy
RUN find . -maxdepth 1 -name "*_test.go" -exec sh -c 'mv "$1" "verified_test_$(basename $1)"' _ {} \; 2>/dev/null || true
COPY verified_test.go ./7361_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
