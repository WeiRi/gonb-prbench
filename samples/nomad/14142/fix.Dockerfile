FROM gonb-nomad-14142-base:latest
ENV GOPROXY=https://goproxy.cn,direct GOSUMDB=off
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /work/upstream
RUN go mod download 2>&1 | tail -3 || true
COPY fix_prod.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff
WORKDIR /work/upstream/nomad
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./14142_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
