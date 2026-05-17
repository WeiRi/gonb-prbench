# base.Dockerfile for etcd-5897 (GOPATH mode)
FROM golang:1.10
ENV GOPATH=/go GO111MODULE=off GO15VENDOREXPERIMENT=1 CGO_ENABLED=1
RUN mkdir -p /go/src/github.com/coreos/etcd
COPY source_full/ /go/src/github.com/coreos/etcd/
COPY fix.diff /tmp/fix.diff
WORKDIR /go/src/github.com/coreos/etcd
# Revert fix.diff to get BUG state
RUN if [ -f /tmp/fix.diff ] && grep -q "^diff --git" /tmp/fix.diff; then \
      git init --quiet && git add -A && \
      git -c user.email=x@x -c user.name=x commit -m base -q && \
      (git apply -R --whitespace=nowarn /tmp/fix.diff 2>&1 || echo "WARNING: fix.diff revert issues"); \
    fi
# Ensure vendor symlinks work
RUN rm -rf cmd/vendor/github.com/coreos/etcd 2>/dev/null; \
    ln -sf cmd/vendor vendor 2>/dev/null; \
    ls vendor/ 2>/dev/null | head -3 || true
