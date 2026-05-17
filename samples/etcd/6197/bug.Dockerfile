FROM gonb-etcd-6197-base-v7:latest
ENV GO111MODULE=off
RUN rm -rf /work 2>/dev/null || true
WORKDIR /go/src/github.com/coreos/etcd/integration
RUN sed -i '/^\/\/ +build cluster_proxy/d' cluster_proxy.go
RUN sed -i '/^\/\/ +build !cluster_proxy/d' cluster_direct.go && rm cluster_direct.go
RUN find . -maxdepth 1 -name "*_test.go" ! -name "verified_test*" ! -name "*race_test*" -delete 2>/dev/null || true
COPY verified_test.go ./6197_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
