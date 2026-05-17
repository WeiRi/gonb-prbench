FROM gonb-kubernetes-134379-base-v2:latest
WORKDIR /work/upstream/pkg/controller/garbagecollector
RUN find . -maxdepth 1 -name "*_test.go" -exec sh -c 'mv "$1" "verified_test_$(basename $1)"' _ {} \; 2>/dev/null || true
COPY verified_test.go ./kubernetes_134379_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
