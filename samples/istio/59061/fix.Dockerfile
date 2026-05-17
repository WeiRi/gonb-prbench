# fix.Dockerfile for istio-59061 (in-place mode)
# Applies the upstream fix (removes `c.Action = action` from Cluster.Run) to the
# base image and re-installs the same in-place regression test. With the racy
# write gone, the test must PASS under -race.
FROM inp-istio-59061

ENV GOTOOLCHAIN=auto \
    GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.org

COPY fix.diff /tmp/fix.diff
RUN cd /go/src/istio.io/istio && git apply --whitespace=nowarn /tmp/fix.diff

COPY verified_test_inplace.go /go/src/istio.io/istio/pkg/kube/multicluster/59061_race_test.go

WORKDIR /go/src/istio.io/istio/pkg/kube/multicluster
