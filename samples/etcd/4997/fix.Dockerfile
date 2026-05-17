# fix.Dockerfile for etcd-4997
FROM gonb-etcd-4997-base
WORKDIR /go/src/github.com/coreos/etcd
COPY fix.diff /tmp/fix.diff
RUN awk 'BEGIN{p=0} /^diff --git/{if ($$0 ~ /consistent_index.go/ && $$0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    (git apply --whitespace=nowarn /tmp/fix_prod.diff 2>/dev/null || \
     (git init --quiet && git add -A && git -c user.email=x@x -c user.name=x commit -m b -q && \
      git apply --whitespace=nowarn /tmp/fix_prod.diff 2>/dev/null) || \
     patch -p1 < /tmp/fix.diff 2>/dev/null || \
     echo "WARNING: fix apply issues")
WORKDIR /go/src/github.com/coreos/etcd/etcdserver
RUN rm -f etcd_4997_race_test.go etcd_4997_race_test.go 2>/dev/null; rm -f *_race_test.go 2>/dev/null; true
RUN for f in *_test.go; do mv "$f" "vt_$$(echo $f)" 2>/dev/null; done; true
COPY verified_test.go ./etcd_4997_race_test.go
RUN go test -race -vet=off -c -o /dev/null . 2>&1 | tail -15 || true
CMD ["sh","-c","go test -race -vet=off -count=10 -timeout=180s -run 'TestRace_4997_ConsistentIndex' ."]
