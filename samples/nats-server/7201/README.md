# nats-server-7201

| Field | Value |
|---|---|
| Project | nats-server |
| Reference | https://github.com/nats-io/nats-server/pull/7201 |
| Bug commit | `3b44c889eef5` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `server/jetstream_api.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000b2038 by goroutine 8:
  ase/nats-server-7201.(*consumer).updatePauseUntil()
      /work/jetstream_api.go:22 +0x44
  ase/nats-server-7201.TestRacePauseUntilRead.func1()
      /work/verified_test.go:32 +0xae
  ase/nats-server-7201.TestRacePauseUntilRead.gowrap1()
      /work/verified_test.go:34 +0x41

Previous read at 0x00c0000b2038 by goroutine 18:
  ase/nats-server-7201.(*consumer).jsConsumerCreateRequest()
      /work/jetstream_api.go:28 +0xe5
  ase/nats-server-7201.TestRacePauseUntilRead.func2()
      /work/verified_test.go:42 +0xdd

Goroutine 8 (running) created at:
  ase/nats-server-7201.TestRacePauseUntilRead()
      /work/verified_test.go:28 +0x17e
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 18 (finished) created at:
  ase/nats-server-7201.TestRacePauseUntilRead()
      /work/verified_test.go:38 +0x2bc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c0000122d0 by goroutine 15:
  ase/nats-server-7201.(*consumer).jsConsumerCreateRequest()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-nats-server-7201-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-nats-server-7201-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-nats-server-7201-fix .
docker run --rm --memory=2g --cpus=1 gonb-nats-server-7201-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-nats-server-7201-bug .
# (then run as above, no --ssh flag)
```
