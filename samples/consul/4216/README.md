# consul-4216

| Field | Value |
|---|---|
| Project | consul |
| Reference | https://github.com/hashicorp/consul/pull/4216 |
| Bug commit | `6e9cbeecd06f` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `agent/consul/client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00019a078 by goroutine 68:
  ase/consul-4216.(*Client).ReloadConfig()
      /work/client.go:31 +0xc4
  ase/consul-4216.TestRPCLimiterRace.func2()
      /work/verified_test.go:27 +0x7a

Previous read at 0x00c00019a078 by goroutine 9:
  ase/consul-4216.(*Client).RPC()
      /work/client.go:26 +0x84
  ase/consul-4216.TestRPCLimiterRace.func1()
      /work/verified_test.go:20 +0x7d

Goroutine 68 (running) created at:
  ase/consul-4216.TestRPCLimiterRace()
      /work/verified_test.go:25 +0xac
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (finished) created at:
  ase/consul-4216.TestRPCLimiterRace()
      /work/verified_test.go:18 +0x1d8
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c0001a81d8 by goroutine 188:
  ase/consul-4216.(*Limiter).Allow()
      /work/client.go:13 +0x97
  ase/consul-4216.(*Client).RPC()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-consul-4216-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-consul-4216-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-consul-4216-fix .
docker run --rm --memory=2g --cpus=1 gonb-consul-4216-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-consul-4216-bug .
# (then run as above, no --ssh flag)
```
