FROM inp-dns-459
COPY fix.diff /tmp/fix.diff
# Filter to apply ONLY client.go (not client_test.go which may be missing/conflicting)
RUN awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /client\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff && \
    cd /go/src/github.com/miekg/dns && \
    git init --quiet 2>/dev/null; git add -A 2>/dev/null; git -c user.email=x@x -c user.name=x commit -m b -q 2>/dev/null; \
    git apply --whitespace=nowarn /tmp/fix_prod.diff && \
    grep -A 3 "shared :=" client.go | head -8
COPY dns_helper.go /go/src/github.com/miekg/dns/dns_helper.go
