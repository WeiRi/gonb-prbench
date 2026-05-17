FROM inp-minio-15770
RUN mkdir -p /work/pr2t-test
COPY go.mod verified_test.go rpc-stats.go /work/pr2t-test/
