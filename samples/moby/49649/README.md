# moby-49649

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/49649 |
| Bug commit | `3e1b15dc9763` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `cmd/docker-proxy/udp_proxy_linux.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00011c390 by goroutine 13:
  runtime.mapaccess1_faststr()
      /usr/local/go/src/runtime/map_faststr.go:13 +0x0
  ase/moby-49649.(*connTrackTable).lookupConn()
      /work/udp_proxy_linux.go:26 +0x108
  ase/moby-49649.TestRace_49649.func2()
      /work/verified_test.go:26 +0xb4

Previous write at 0x00c00011c390 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/moby-49649.(*connTrackTable).addConn()
      /work/udp_proxy_linux.go:20 +0xf7
  ase/moby-49649.TestRace_49649.func1()
      /work/verified_test.go:20 +0xb1

Goroutine 13 (running) created at:
  ase/moby-49649.TestRace_49649()
      /work/verified_test.go:23 +0xad
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1446 +0x216
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1493 +0x47

Goroutine 8 (finished) created at:
  ase/moby-49649.TestRace_49649()
      /work/verified_test.go:17 +0x176
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1446 +0x216
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1493 +0x47
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-49649-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-49649-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-49649-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-49649-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-49649-bug .
# (then run as above, no --ssh flag)
```
