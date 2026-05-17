FROM gonb-consul-2262-base:latest
RUN rm -rf /work 2>/dev/null || true
WORKDIR /go/src/github.com/hashicorp/consul/command/agent
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./2262_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
