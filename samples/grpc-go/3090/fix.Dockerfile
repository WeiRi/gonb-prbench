FROM gonb-grpc-go-3090-base:latest
RUN rm -rf /work 2>/dev/null || true
WORKDIR /go/src/google.golang.org/grpc
COPY fix_prod.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test_fixed.go ./3090_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
