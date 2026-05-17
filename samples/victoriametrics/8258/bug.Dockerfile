FROM gonb-victoriametrics-8258-base:latest
WORKDIR /work/upstream/app/vmalert/notifier
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null
COPY verified_test.go ./vm_8258_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
