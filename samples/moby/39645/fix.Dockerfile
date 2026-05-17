# fix.Dockerfile for moby-39645 — apply fix.diff on top of BUG worktree
FROM gonb-moby-39645-bug:latest

ENV GOPATH=/go GO111MODULE=off

# Remove old pr2t-test to prevent gate from picking up stub tests
RUN rm -rf /work/pr2t-test

# Set up GOPATH with source in correct location for vendor resolution
RUN mkdir -p /go/src/github.com/docker && \
    cp -r /work/upstream /go/src/github.com/docker/docker

# Apply fix.diff to get FIX state
WORKDIR /go/src/github.com/docker/docker
COPY fix.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff || patch -p1 < /tmp/fix.diff

WORKDIR /go/src/github.com/docker/docker/container

RUN find . -maxdepth 1 -name "*_test.go" \
    ! -name "*_race_test*" ! -name "*inplace_test*" ! -name "verified_test*" \
    -delete 2>/dev/null || true && rm -f race_test.go

COPY verified_test.go ./moby_39645_race_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
