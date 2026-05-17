FROM gonb-kubernetes-133781-base-v3:latest
RUN rm -rf /work/pr2t-test 2>/dev/null || true
WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /helpers\.go|priorityclass\/storage\/storage\.go|defaulttolerationseconds\/admission\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    patch -p1 -i /tmp/fix_prod.diff
WORKDIR /work/upstream/pkg/apis/scheduling/v1
RUN find . -maxdepth 1 -name "*_test.go" -exec sh -c 'mv "$1" "verified_test_$(basename $1)"' _ {} \; 2>/dev/null || true
COPY verified_test.go ./133781_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
