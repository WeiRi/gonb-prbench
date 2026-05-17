# bug.Dockerfile for nomad-14121 — preserve upstream's own deploymentwatcher tests
FROM gonb-nomad-14121-base:latest

WORKDIR /work/upstream/nomad/deploymentwatcher

RUN mv deployments_watcher_test.go verified_test_dw_test.go && \
    mv testutil_test.go verified_test_dwutil_test.go && \
    ls *_test.go | head

RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
