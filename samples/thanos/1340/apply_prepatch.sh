#!/bin/sh
# Apply thanos-1340 fix from fix.diff, then prepatch to move mutex earlier in Put()
set -e
cd /work/upstream

# Apply original fix
awk 'BEGIN{p=0} /^diff --git/{if ($0 ~ /pool\.go/ && $0 !~ /_test\.go/) p=1; else p=0} p==1' /tmp/fix.diff > /tmp/fix_prod.diff
(git apply --whitespace=nowarn /tmp/fix_prod.diff || (git init --quiet && git add -A && git -c user.email=x@x -c user.name=x commit -m b -q && git apply --whitespace=nowarn /tmp/fix_prod.diff))

# Pre-patch: move p.mtx.Lock()/defer p.mtx.Unlock() from bottom to top of Put()
# Line numbers in FIXED pool.go:
#   98: p.mtx.Lock()
#   99: defer p.mtx.Unlock()
#   84: func (p *BytesPool) Put(b *[]byte) {
cd /work/upstream/pkg/pool
TAB="$(printf '\t')"
# Add mutex after line 84 (function opening)
sed -i "84a\\${TAB}p.mtx.Lock()\\n${TAB}defer p.mtx.Unlock()" pool.go
# Remove lines 100,101 (shifted from 98,99 after adding 2 lines)
sed -i '100,101d' pool.go
