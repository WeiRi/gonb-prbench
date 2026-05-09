# syncthing-8321

| Field | Value |
|---|---|
| Project | syncthing |
| Reference | https://github.com/syncthing/syncthing/pull/8321 |
| Bug commit | `31a78592e80e` |
| Category | anonymous_function |
| Oracle | RACE |
| Primary diff file | `lib/connections/service.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000216000 by goroutine 10:
  ase/syncthing-8321.(*service).HandleConns.func1()
      /work/service.go:34 +0x9e

Previous write at 0x00c000216000 by goroutine 9:
  ase/syncthing-8321.(*service).HandleConns()
      /work/service.go:30 +0x155
  ase/syncthing-8321.TestRace_8321.func2()
      /work/verified_test.go:33 +0xa4

Goroutine 10 (running) created at:
  ase/syncthing-8321.(*service).HandleConns()
      /work/service.go:32 +0x73
  ase/syncthing-8321.TestRace_8321.func2()
      /work/verified_test.go:33 +0xa4

Goroutine 9 (running) created at:
  ase/syncthing-8321.TestRace_8321()
      /work/verified_test.go:31 +0x356
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c000216000 by goroutine 9:
  ase/syncthing-8321.(*service).HandleConns()
      /work/service.go:30 +0x155
  ase/syncthing-8321.TestRace_8321.func2()
      /work/verified_test.go:33 +0xa4

Previous read at 0x00c000216000 by goroutine 12:
  ase/syncthing-8321.(*service).HandleConns.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-syncthing-8321-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-syncthing-8321-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-syncthing-8321-fix .
docker run --rm --memory=2g --cpus=1 gonb-syncthing-8321-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-syncthing-8321-bug .
# (then run as above, no --ssh flag)
```
