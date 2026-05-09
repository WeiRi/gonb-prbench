# syncthing-4526

| Field | Value |
|---|---|
| Project | syncthing |
| Reference | https://github.com/syncthing/syncthing/pull/4526 |
| Bug commit | `9471b9f6af02` |
| Category | anonymous_function |
| Oracle | RACE |
| Primary diff file | `lib/connections/service.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00020c000 by goroutine 55:
  ase/syncthing-4526.DialParallel.func1()
      /work/service.go:29 +0x45

Previous write at 0x00c00020c000 by goroutine 8:
  ase/syncthing-4526.DialParallel()
      /work/service.go:26 +0xf9
  ase/syncthing-4526.TestRace_4526.func1()
      /work/verified_test.go:18 +0xce

Goroutine 55 (running) created at:
  ase/syncthing-4526.DialParallel()
      /work/service.go:28 +0x9e
  ase/syncthing-4526.TestRace_4526.func1()
      /work/verified_test.go:18 +0xce

Goroutine 8 (running) created at:
  ase/syncthing-4526.TestRace_4526()
      /work/verified_test.go:16 +0x14d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_4526 (0.01s)
FAIL
FAIL	ase/syncthing-4526	0.034s
FAIL
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-syncthing-4526-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-syncthing-4526-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-syncthing-4526-fix .
docker run --rm --memory=2g --cpus=1 gonb-syncthing-4526-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-syncthing-4526-bug .
```
