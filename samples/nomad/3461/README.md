# nomad-3461

| Field | Value |
|---|---|
| Project | nomad |
| Reference | https://github.com/hashicorp/nomad/pull/3461 |
| Bug commit | `2464b02aa2a1` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `client/client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d4570 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/nomad-3461.(*Client).updateAttributes()
      /work/client.go:34 +0xa4
  ase/nomad-3461.TestRaceNodeAttributes.func1()
      /work/verified_test.go:39 +0xb8
  ase/nomad-3461.TestRaceNodeAttributes.gowrap1()
      /work/verified_test.go:41 +0x41

Previous read at 0x00c0000d4570 by goroutine 13:
  ase/nomad-3461.TestRaceNodeAttributes.func2()
      /work/verified_test.go:49 +0x109

Goroutine 8 (running) created at:
  ase/nomad-3461.TestRaceNodeAttributes()
      /work/verified_test.go:36 +0x3e4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 13 (running) created at:
  ase/nomad-3461.TestRaceNodeAttributes()
      /work/verified_test.go:45 +0x534
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c00020c098 by goroutine 8:
  ase/nomad-3461.(*Client).updateAttributes()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-nomad-3461-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-nomad-3461-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-nomad-3461-fix .
docker run --rm --memory=2g --cpus=1 gonb-nomad-3461-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-nomad-3461-bug .
# (then run as above, no --ssh flag)
```
