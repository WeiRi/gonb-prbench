# syncthing-6477

| Field | Value |
|---|---|
| Project | syncthing |
| Reference | https://github.com/syncthing/syncthing/pull/6477 |
| Bug commit | `7709ac33a7fe` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `lib/util/utils.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000020310 by goroutine 8:
  ase/syncthing-6477.(*service).Serve()
      /work/lib_util_utils.go:22 +0x59
  ase/syncthing-6477.TestRaceStoppedField.func1()
      /work/verified_test.go:14 +0x84

Previous read at 0x00c000020310 by goroutine 13:
  ase/syncthing-6477.(*service).Stop()
      /work/lib_util_utils.go:35 +0x48
  ase/syncthing-6477.TestRaceStoppedField.func2()
      /work/verified_test.go:15 +0x84

Goroutine 8 (running) created at:
  ase/syncthing-6477.TestRaceStoppedField()
      /work/verified_test.go:14 +0x1d0
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 13 (finished) created at:
  ase/syncthing-6477.TestRaceStoppedField()
      /work/verified_test.go:15 +0xfe
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: TestRaceStoppedField (0.01s)
    testing.go:1398: race detected during execution of test
FAIL
FAIL	ase/syncthing-6477	0.032s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-syncthing-6477-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-syncthing-6477-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-syncthing-6477-fix .
docker run --rm --memory=2g --cpus=1 gonb-syncthing-6477-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-syncthing-6477-bug .
# (then run as above, no --ssh flag)
```
