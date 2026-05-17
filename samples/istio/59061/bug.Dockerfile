# bug.Dockerfile for istio-59061 (in-place mode)
# Base image already has upstream source at bug commit f6dcc1843963 with
# `c.Action = action` present at pkg/kube/multicluster/cluster.go:99.
# We overwrite the placeholder test with the real in-place test that
# concurrently invokes Cluster.Run() so the race detector reports a frame in
# cluster.go (the exact line removed by the PR fix).
FROM inp-istio-59061

ENV GOTOOLCHAIN=auto \
    GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.org

# Install the real in-place regression test.
COPY verified_test_inplace.go /go/src/istio.io/istio/pkg/kube/multicluster/59061_race_test.go

WORKDIR /go/src/istio.io/istio/pkg/kube/multicluster
