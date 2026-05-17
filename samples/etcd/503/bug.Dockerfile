FROM gonb-etcd-503-base-v7:latest
ENV GO111MODULE=off
WORKDIR /go/src/github.com/coreos/etcd/mod/lock/v2
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./503_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
