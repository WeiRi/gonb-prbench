FROM inp-minio-15770
COPY fix.diff /tmp/fix.diff
RUN cd /go/src/github.com/minio/minio && git apply --whitespace=nowarn /tmp/fix.diff 2>/dev/null || true
RUN mkdir -p /work/pr2t-test
COPY go.mod verified_test.go /work/pr2t-test/
# Use the FIXED version of rpc-stats.go (with atomic operations)
COPY rpc-stats-fixed.go /work/pr2t-test/rpc-stats.go
