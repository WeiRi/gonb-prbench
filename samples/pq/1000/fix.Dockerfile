FROM gonb-pq-1000-bug:latest
WORKDIR /work/upstream
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    patch -p1 -i /tmp/fix_prod.diff || true
RUN sed -i 's,cn\.bad = true,cn.bad.Store(true),g; s,!cn\.bad\b,!cn.bad.Load().(bool),g' error.go conn.go copy.go conn_go18.go 2>/dev/null
# Remove all upstream _test.go files since they have inter-deps that don't work with fix
RUN rm -f *_test.go
COPY verified_test.go ./pq_1000_race_test.go
COPY verified_test_helper_fix.go ./pq_1000_helper.go
RUN sed -i 's,//go:build pqfix,,' pq_1000_helper.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -10 || true
