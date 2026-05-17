FROM gonb-hugo-7393-base-v2:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /work/upstream
COPY fix_prod.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff
WORKDIR /work/upstream/hugolib
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test_fixed.go ./7393_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
