# bug.Dockerfile for moby-26695 (GOPATH-era, vendor mode)
FROM gonb-moby-26695-base:latest

ENV GOPATH=/go/src/github.com/docker/docker/vendor:/go
ENV GO15VENDOREXPERIMENT=1

# Revert fix.diff to get BUG state
WORKDIR /go/src/github.com/docker/docker
COPY fix.diff /tmp/fix.diff
RUN git apply -R --whitespace=nowarn /tmp/fix.diff 2>/dev/null || patch -p1 -R < /tmp/fix.diff 2>/dev/null || true

WORKDIR /go/src/github.com/docker/docker/libcontainerd
# Delete all _test.go (including empty leftover race_test.go), then plant ours
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true

COPY verified_test.go ./moby_26695_race_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
