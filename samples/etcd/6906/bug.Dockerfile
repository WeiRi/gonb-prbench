# bug.Dockerfile for etcd-6906
FROM gonb-etcd-6906-base
WORKDIR /go/src/github.com/coreos/etcd/proxy/grpcproxy
RUN rm -f etcd_6906_race_test.go etcd_6906_race_test.go 2>/dev/null; rm -f *_race_test.go 2>/dev/null; true
COPY verified_test.go ./etcd_6906_race_test.go
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=180s -run 'TestRace_6906' ."]
