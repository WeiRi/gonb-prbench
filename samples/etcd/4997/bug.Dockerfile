# bug.Dockerfile for etcd-4997
FROM gonb-etcd-4997-base
WORKDIR /go/src/github.com/coreos/etcd/etcdserver
RUN rm -f etcd_4997_race_test.go etcd_4997_race_test.go 2>/dev/null; rm -f *_race_test.go 2>/dev/null; true
COPY verified_test.go ./etcd_4997_race_test.go
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=180s -run 'TestRace_4997_ConsistentIndex' ."]
