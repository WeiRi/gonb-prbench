# syncthing-8042

| Field | Value |
|---|---|
| Project | syncthing |
| Reference | https://github.com/syncthing/syncthing/pull/8042 |
| Bug commit | `dec6f80d2bc9` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `cmd/strelaysrv/main.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000bc8e0 by goroutine 8:
  ase/syncthing-8042.(*Service).process()
      /work/service.go:33 +0x12d
  ase/syncthing-8042.TestMappingRace.func1()
      /work/verified_test.go:26 +0x99

Previous write at 0x00c0000bc8e0 by goroutine 9:
  ase/syncthing-8042.(*Service).updateMapping()
      /work/service.go:42 +0x64
  ase/syncthing-8042.TestMappingRace.func2()
      /work/verified_test.go:30 +0xa6

Goroutine 8 (running) created at:
  ase/syncthing-8042.TestMappingRace()
      /work/verified_test.go:24 +0x32c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (finished) created at:
  ase/syncthing-8042.TestMappingRace()
      /work/verified_test.go:28 +0x2b
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c0000d4540 by goroutine 8:
  runtime.mapiterinit()
      /usr/local/go/src/runtime/map.go:816 +0x0
  ase/syncthing-8042.(*Service).process()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-syncthing-8042-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-syncthing-8042-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-syncthing-8042-fix .
docker run --rm --memory=2g --cpus=1 gonb-syncthing-8042-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-syncthing-8042-bug .
# (then run as above, no --ssh flag)
```
