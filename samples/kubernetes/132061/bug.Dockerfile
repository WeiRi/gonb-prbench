FROM gonb-kubernetes-132061-base-v2:latest
ENV GOWORK=off
WORKDIR /work/upstream/staging/src/k8s.io/apiserver/pkg/cel/common
RUN find . -maxdepth 1 -name "*_test.go" -exec sh -c 'mv "$1" "verified_test_$(basename $1)"' _ {} \; 2>/dev/null || true
COPY verified_test.go ./kubernetes_132061_race_test.go
RUN GOWORK=off go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
