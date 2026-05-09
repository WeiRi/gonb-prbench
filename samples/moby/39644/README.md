# moby-39644

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/39644 |
| Bug commit | `928381b2215c` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `daemon/graphdriver/quota/projectquota.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000a2568 by goroutine 14:
  ase/moby-39644.(*Control).SetQuota()
      /work/projectquota.go:7 +0xcc
  ase/moby-39644.TestRace_39644.func1()
      /work/verified_test.go:21 +0xb8

Previous write at 0x00c0000a2568 by goroutine 8:
  ase/moby-39644.(*Control).SetQuota()
      /work/projectquota.go:8 +0x144
  ase/moby-39644.TestRace_39644.func1()
      /work/verified_test.go:21 +0xb8

Goroutine 14 (running) created at:
  ase/moby-39644.TestRace_39644()
      /work/verified_test.go:18 +0x196
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1446 +0x216
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1493 +0x47

Goroutine 8 (finished) created at:
  ase/moby-39644.TestRace_39644()
      /work/verified_test.go:18 +0x196
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1446 +0x216
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1493 +0x47
==================
==================
WARNING: DATA RACE
Read at 0x00c00010e088 by goroutine 9:
  ase/moby-39644.(*Control).GetQuota()
      /work/projectquota.go:14 +0xdd
  ase/moby-39644.TestRace_39644.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-39644-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-39644-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-39644-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-39644-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-39644-bug .
# (then run as above, no --ssh flag)
```
