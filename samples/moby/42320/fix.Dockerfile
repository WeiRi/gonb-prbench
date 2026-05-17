FROM gonb-moby-42320-base-v6:latest
RUN rm -rf /work 2>/dev/null || true
WORKDIR /go/src/github.com/docker/docker
COPY fix.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff
WORKDIR /go/src/github.com/docker/docker/daemon/logger
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./42320_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
