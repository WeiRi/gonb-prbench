# moby-22966

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/22966 |
| Bug commit | `29dbcbad8784` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/discovery/memory/memory.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00019c120 by goroutine 10:
  ase/moby-22966.(*Discovery).Register()
      /work/memory.go:19 +0x1b1
  ase/moby-22966.TestRace_22966.func1()
      /work/verified_test.go:21 +0x109
  ase/moby-22966.TestRace_22966.gowrap1()
      /work/verified_test.go:23 +0x41

Previous write at 0x00c00019c120 by goroutine 8:
  ase/moby-22966.(*Discovery).Register()
      /work/memory.go:19 +0x268
  ase/moby-22966.TestRace_22966.func1()
      /work/verified_test.go:21 +0x109
  ase/moby-22966.TestRace_22966.gowrap1()
      /work/verified_test.go:23 +0x41

Goroutine 10 (running) created at:
  ase/moby-22966.TestRace_22966()
      /work/verified_test.go:18 +0xb9
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/moby-22966.TestRace_22966()
      /work/verified_test.go:18 +0xb9
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-22966-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-22966-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-22966-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-22966-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-22966-bug .
# (then run as above, no --ssh flag)
```
