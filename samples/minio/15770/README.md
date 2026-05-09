# minio-15770

| Field | Value |
|---|---|
| Project | minio |
| Reference | https://github.com/minio/minio/pull/15770 |
| Bug commit | `d44f3526dc71` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `internal/rest/rpc-stats.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000121c8 by goroutine 35:
  buggy.SetupReqStatsUpdate.func3()
      /work/rpc_stats.go:41 +0x47
  buggy.TestRaceTCPDialStats.func3()
      /work/race_test.go:64 +0x1a7

Previous write at 0x00c0000121c8 by goroutine 165:
  buggy.SetupReqStatsUpdate.func1()
      /work/rpc_stats.go:31 +0x47
  net.(*sysDialer).dialSingle()
      /usr/local/go/src/net/dial.go:636 +0x16c
  net.(*sysDialer).dialSerial()
      /usr/local/go/src/net/dial.go:614 +0x291
  net.(*sysDialer).dialParallel()
      /usr/local/go/src/net/dial.go:515 +0x5da
  net.(*Dialer).DialContext()
      /usr/local/go/src/net/dial.go:506 +0xb35
  net/http.(*Transport).dial()
      /usr/local/go/src/net/http/transport.go:1196 +0x2d6
  net/http.(*Transport).dialConn()
      /usr/local/go/src/net/http/transport.go:1625 +0xdc4
  net/http.(*Transport).dialConnFor()
      /usr/local/go/src/net/http/transport.go:1467 +0x129
  net/http.(*Transport).queueForDial.func1()
      /usr/local/go/src/net/http/transport.go:1436 +0x44

Goroutine 35 (running) created at:
  buggy.TestRaceTCPDialStats()
      /work/race_test.go:52 +0x5bb
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-minio-15770-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-minio-15770-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-minio-15770-fix .
docker run --rm --memory=2g --cpus=1 gonb-minio-15770-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-minio-15770-bug .
# (then run as above, no --ssh flag)
```
