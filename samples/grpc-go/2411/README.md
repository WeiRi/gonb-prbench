# grpc-go-2411

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/2411 |
| Bug commit | `0430365f23cf` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `internal/channelz/funcs.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000190420 by goroutine 14:
  ase/grpc-go-2411.(*ChannelMap).DeleteSelfFromMap()
      /work/funcs.go:51 +0x89
  ase/grpc-go-2411.TestRace_PR2411_ChannelzGetChannel.func2()
      /work/verified_test.go:66 +0xdb

Previous read at 0x00c000190420 by goroutine 9:
  ase/grpc-go-2411.(*ChannelMap).GetChannel()
      /work/funcs.go:45 +0x8e
  ase/grpc-go-2411.TestRace_PR2411_ChannelzGetChannel.func1()
      /work/verified_test.go:58 +0xbd

Goroutine 14 (running) created at:
  ase/grpc-go-2411.TestRace_PR2411_ChannelzGetChannel()
      /work/verified_test.go:63 +0x2d0
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (running) created at:
  ase/grpc-go-2411.TestRace_PR2411_ChannelzGetChannel()
      /work/verified_test.go:55 +0x144
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR2411_ChannelzGetChannel (0.08s)
FAIL
FAIL	ase/grpc-go-2411	0.096s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-grpc-go-2411-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-2411-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-grpc-go-2411-fix .
docker run --rm --memory=2g --cpus=1 gonb-grpc-go-2411-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-grpc-go-2411-bug .
# (then run as above, no --ssh flag)
```
