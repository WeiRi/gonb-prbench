FROM gonb-go-ethereum-14940-bug-inplace:latest
RUN rm -rf /work 2>/dev/null || true
WORKDIR /go/src/github.com/ethereum/go-ethereum
COPY prepatch.diff /tmp/prepatch.diff
COPY fix_prod.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/prepatch.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff
WORKDIR /go/src/github.com/ethereum/go-ethereum/core
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./14940_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
