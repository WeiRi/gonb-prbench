# teleport-65910

| Field | Value |
|---|---|
| Project | teleport |
| Reference | https://github.com/gravitational/teleport/pull/65910 |
| Bug commit | `1bd317b98340` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `lib/srv/app/gcp/handler.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00028ed78 by goroutine 73:
  github.com/gravitational/teleport/lib/srv/app/gcp.(*handler).getToken.func1()
      /workspace/lib/srv/app/gcp/handler.go:274 +0x244

Previous read at 0x00c00028ed78 by goroutine 52:
  github.com/gravitational/teleport/lib/srv/app/gcp.(*handler).getToken()
      /workspace/lib/srv/app/gcp/handler.go:284 +0x527
  github.com/gravitational/teleport/lib/srv/app/gcp.TestRace_65910.func2()
      /workspace/lib/srv/app/gcp/handler_race_test.go:87 +0x239
  github.com/gravitational/teleport/lib/srv/app/gcp.TestRace_65910.gowrap1()
      /workspace/lib/srv/app/gcp/handler_race_test.go:89 +0x41

Goroutine 73 (running) created at:
  github.com/gravitational/teleport/lib/srv/app/gcp.(*handler).getToken()
      /workspace/lib/srv/app/gcp/handler.go:273 +0x3aa
  github.com/gravitational/teleport/lib/srv/app/gcp.TestRace_65910.func2()
      /workspace/lib/srv/app/gcp/handler_race_test.go:87 +0x239
  github.com/gravitational/teleport/lib/srv/app/gcp.TestRace_65910.gowrap1()
      /workspace/lib/srv/app/gcp/handler_race_test.go:89 +0x41

Goroutine 52 (running) created at:
  github.com/gravitational/teleport/lib/srv/app/gcp.TestRace_65910()
      /workspace/lib/srv/app/gcp/handler_race_test.go:80 +0x285
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c000a0f3c0 by goroutine 73:
  github.com/gravitational/teleport/lib/srv/app/gcp.(*handler).getToken.func1()
      /workspace/lib/srv/app/gcp/handler.go:274 +0x288
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-teleport-65910-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-teleport-65910-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-teleport-65910-fix .
docker run --rm --memory=2g --cpus=1 gonb-teleport-65910-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-teleport-65910-bug .
# (then run as above, no --ssh flag)
```
