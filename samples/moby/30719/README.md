# moby-30719

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/30719 |
| Bug commit | `c3b660b11280` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `daemon/graphdriver/counter.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000208008 by goroutine 38:
  ase/moby-30719.(*RefCounter).Decrement()
      /work/counter.go:47 +0x124
  ase/moby-30719.TestRace_30719.func2()
      /work/verified_test.go:30 +0x165

Previous write at 0x00c000208008 by goroutine 13:
  ase/moby-30719.(*RefCounter).Increment()
      /work/counter.go:35 +0x148
  ase/moby-30719.TestRace_30719.func1()
      /work/verified_test.go:22 +0xa5

Goroutine 38 (running) created at:
  ase/moby-30719.TestRace_30719()
      /work/verified_test.go:27 +0x211
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 13 (running) created at:
  ase/moby-30719.TestRace_30719()
      /work/verified_test.go:19 +0x145
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c000208008 by goroutine 61:
  ase/moby-30719.(*RefCounter).Decrement()
      /work/counter.go:47 +0x124
  ase/moby-30719.TestRace_30719.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-30719-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-30719-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-30719-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-30719-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-30719-bug .
# (then run as above, no --ssh flag)
```
