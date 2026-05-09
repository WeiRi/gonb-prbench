# syncthing-10430

| Field | Value |
|---|---|
| Project | syncthing |
| Reference | https://github.com/syncthing/syncthing/pull/10430 |
| Bug commit | `20d2406a0eb9` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `lib/fs/casefs.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000d4540 by goroutine 8:
  runtime.mapaccess2_faststr()
      /usr/local/go/src/runtime/map_faststr.go:108 +0x0
  ase/syncthing-10430.(*defaultRealCaser).getExpireAdd()
      /work/casefs.go:29 +0x11b
  ase/syncthing-10430.TestCaseCacheRace.func1()
      /work/verified_test.go:23 +0xa5

Previous write at 0x00c0000d4540 by goroutine 9:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/syncthing-10430.(*defaultRealCaser).getExpireAdd()
      /work/casefs.go:33 +0x1b4
  ase/syncthing-10430.TestCaseCacheRace.func2()
      /work/verified_test.go:27 +0xa5

Goroutine 8 (running) created at:
  ase/syncthing-10430.TestCaseCacheRace()
      /work/verified_test.go:21 +0x29c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (finished) created at:
  ase/syncthing-10430.TestCaseCacheRace()
      /work/verified_test.go:25 +0x18d
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-syncthing-10430-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-syncthing-10430-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-syncthing-10430-fix .
docker run --rm --memory=2g --cpus=1 gonb-syncthing-10430-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-syncthing-10430-bug .
# (then run as above, no --ssh flag)
```
