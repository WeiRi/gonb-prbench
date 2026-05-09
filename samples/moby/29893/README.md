# moby-29893

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/29893 |
| Bug commit | `9c3955aae1bf` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/plugins/plugins.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
# WHITE_BOX_UPSTREAM_NOT_REPRODUCIBLE
# Sample: moby-29893
# Reason: race_detector_timing_window_missed_in_container_300s
# Detail: docker re-run on upstream failed; v3 worker attempted with --memory=12g --cpus=6 --timeout=900s
# Stderr: 
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-29893-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-29893-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-29893-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-29893-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-29893-bug .
# (then run as above, no --ssh flag)
```
