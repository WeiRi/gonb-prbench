FROM gonb-kubernetes-100638-base-v3:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /work/upstream
RUN go mod download 2>&1 | tail -5 || true
COPY fix.diff /tmp/fix.diff
RUN git apply --whitespace=nowarn /tmp/fix.diff
WORKDIR /work/upstream/staging/src/k8s.io/apiserver/pkg/util/flowcontrol/fairqueuing/queueset
RUN find . -maxdepth 1 -name "*_test.go" -delete 2>/dev/null || true
COPY verified_test.go ./100638_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
