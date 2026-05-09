# nats-server-5146

| Field | Value |
|---|---|
| Project | nats-server |
| Reference | https://github.com/nats-io/nats-server/pull/5146 |
| Bug commit | `d9ded28ba4ec` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `server/consumer.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00012a1d8 by goroutine 103:
  ase/nats-server-5146.TestRaceConsumerStopWithFlagsAccJs.func2()
      /work/verified_test.go:70 +0xfc

Previous read at 0x00c00012a1d8 by goroutine 41:
  ase/nats-server-5146.TestRaceConsumerStopWithFlagsAccJs.func1()
      /work/verified_test.go:53 +0xf5

Goroutine 103 (running) created at:
  ase/nats-server-5146.TestRaceConsumerStopWithFlagsAccJs()
      /work/verified_test.go:65 +0x2c4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44

Goroutine 41 (running) created at:
  ase/nats-server-5146.TestRaceConsumerStopWithFlagsAccJs()
      /work/verified_test.go:45 +0x190
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c00012a1d8 by goroutine 103:
  ase/nats-server-5146.TestRaceConsumerStopWithFlagsAccJs.func2()
      /work/verified_test.go:73 +0x164

Previous read at 0x00c00012a1d8 by goroutine 18:
  ase/nats-server-5146.TestRaceConsumerStopWithFlagsAccJs.func1()
      /work/verified_test.go:53 +0xf5
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-nats-server-5146-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-nats-server-5146-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-nats-server-5146-fix .
docker run --rm --memory=2g --cpus=1 gonb-nats-server-5146-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-nats-server-5146-bug .
# (then run as above, no --ssh flag)
```
