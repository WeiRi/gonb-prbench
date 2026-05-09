# moby-28290

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/28290 |
| Bug commit | `956ff8f773d3` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `volume/store/store.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000126420 by goroutine 15:
  runtime.mapaccess1_faststr()
      /usr/local/go/src/runtime/map_faststr.go:13 +0x0
  ase/moby-28290.(*_BugVolumeStore).list()
      /work/verified_test.go:49 +0xf9
  ase/moby-28290.TestRace_PR28290_VolumeStoreMapAccess.func2()
      /work/verified_test.go:104 +0xd0

Previous write at 0x00c000126420 by goroutine 11:
  runtime.mapdelete_faststr()
      /usr/local/go/src/runtime/map_faststr.go:301 +0x0
  ase/moby-28290.(*_BugVolumeStore).removeVolume()
      /work/verified_test.go:63 +0xa4
  ase/moby-28290.TestRace_PR28290_VolumeStoreMapAccess.func1()
      /work/verified_test.go:89 +0xbd

Goroutine 15 (running) created at:
  ase/moby-28290.TestRace_PR28290_VolumeStoreMapAccess()
      /work/verified_test.go:96 +0x2f7
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 11 (running) created at:
  ase/moby-28290.TestRace_PR28290_VolumeStoreMapAccess()
      /work/verified_test.go:81 +0x1d7
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-28290-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-28290-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-28290-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-28290-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-28290-bug .
# (then run as above, no --ssh flag)
```
