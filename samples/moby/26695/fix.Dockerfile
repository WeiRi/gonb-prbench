# fix.Dockerfile for moby-26695 (worktree already at FIX state, just wire test)
FROM gonb-moby-26695-base:latest

ENV GOPATH=/go/src/github.com/docker/docker/vendor:/go
ENV GO15VENDOREXPERIMENT=1

WORKDIR /go/src/github.com/docker/docker/libcontainerd
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true

COPY verified_test.go ./moby_26695_race_test.go

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
