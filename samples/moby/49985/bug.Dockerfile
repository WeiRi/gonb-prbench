FROM gonb-moby-49985-base-v8:latest
RUN rm -rf /work 2>/dev/null || true
WORKDIR /go/src/github.com/docker/docker/libnetwork/networkdb
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./49985_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
