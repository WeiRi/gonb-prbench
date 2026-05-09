# moby-21677

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/21677 |
| Bug commit | `81d9eaa27e4e` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `layer/ro_layer.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000012148 by goroutine 14:
  ase/moby-21677.(*roLayer).hold()
      /work/ro_layer.go:29 +0xac
  ase/moby-21677.TestRace_21677.func1()
      /work/verified_test.go:21 +0xa4

Previous write at 0x00c000012148 by goroutine 8:
  ase/moby-21677.(*roLayer).hold()
      /work/ro_layer.go:29 +0xc4
  ase/moby-21677.TestRace_21677.func1()
      /work/verified_test.go:21 +0xa4

Goroutine 14 (running) created at:
  ase/moby-21677.TestRace_21677()
      /work/verified_test.go:18 +0x245
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/moby-21677.TestRace_21677()
      /work/verified_test.go:18 +0x245
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c000012148 by goroutine 14:
  ase/moby-21677.(*roLayer).hold()
      /work/ro_layer.go:29 +0xc4
  ase/moby-21677.TestRace_21677.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-21677-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-21677-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-21677-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-21677-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-21677-bug .
# (then run as above, no --ssh flag)
```
