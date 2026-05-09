# kubernetes-132049

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/132049 |
| Bug commit | `0f478e52a8da` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/apiserver/pkg/endpoints/handlers/patch.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
# WHITE_BOX_UPSTREAM_NOT_REPRODUCIBLE
# Sample: kubernetes-132049
# Reason: upstream_go_mod_broken_in_pre_modules
# Detail: docker re-run on upstream failed; v3 worker attempted with --memory=12g --cpus=6 --timeout=900s
# Stderr: go: errors parsing go.mod:|go: errors parsing ../../../go.mod:|
# 
# This sample's race trace was originally captured during dataset construction
# (see marker.real_frame_hits + marker.bug_races). The current dataset<sid>/
# verified_test.go references upstream types and cannot trigger the race in standalone
# docker without the full upstream source tree. The marker remains ab_class=A based on
# the original docker validation.
# 
# Mitigation: in the public artifact's REPRO_GUIDE.md, note that ~10% of A samples
# (24 / 263) require an upstream checkout. The remaining ~90% (239 / 263) reproduce
# standalone via samples/<sid>/run.sh.
```

(Full trace in `race_report_bug.txt`.)

## How to reproduce

### 1. SSH agent setup (one-time)
```bash
eval $(ssh-agent -a /tmp/ssh-agent-gonb.sock)
ssh-add ~/.ssh/id_ed25519
export SSH_AUTH_SOCK=/tmp/ssh-agent-gonb.sock
```

### 2. Build bug image
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-132049-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-132049-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-132049-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-132049-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-132049-bug .
# (then run as above, no --ssh flag)
```
