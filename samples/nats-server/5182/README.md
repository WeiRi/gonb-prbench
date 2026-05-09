# nats-server-5182

| Field | Value |
|---|---|
| Project | nats-server |
| Reference | https://github.com/nats-io/nats-server/pull/5182 |
| Bug commit | `d368a6b1f1f3` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `server/jetstream_cluster.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000f8050 by goroutine 9:
  ase/nats-server-5182.(*stream).updateLseq()
      /work/jetstream_cluster.go:17 +0x56
  ase/nats-server-5182.TestRaceLseqRead.func1()
      /work/verified_test.go:29 +0x99

Previous read at 0x00c0000f8050 by goroutine 19:
  ase/nats-server-5182.(*stream).processClusteredInboundMsg()
      /work/jetstream_cluster.go:25 +0x99
  ase/nats-server-5182.TestRaceLseqRead.func2()
      /work/verified_test.go:38 +0x9e

Goroutine 9 (running) created at:
  ase/nats-server-5182.TestRaceLseqRead()
      /work/verified_test.go:26 +0x84
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 19 (running) created at:
  ase/nats-server-5182.TestRaceLseqRead()
      /work/verified_test.go:35 +0x16b
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRaceLseqRead (0.10s)
FAIL
FAIL	ase/nats-server-5182	0.131s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-nats-server-5182-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-nats-server-5182-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-nats-server-5182-fix .
docker run --rm --memory=2g --cpus=1 gonb-nats-server-5182-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-nats-server-5182-bug .
# (then run as above, no --ssh flag)
```
