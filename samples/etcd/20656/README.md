# etcd-20656

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/20656 |
| Bug commit | `4b16856c96c7` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `tests/robustness/client/client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00019e550 by goroutine 14:
  ase/etcd-20656.(*RecordingClient).AppendRespUnsafe()
      /work/client.go:18 +0xdd
  ase/etcd-20656.TestRace_20656.func1()
      /work/verified_test.go:24 +0xad
  ase/etcd-20656.TestRace_20656.gowrap1()
      /work/verified_test.go:28 +0x41

Previous write at 0x00c00019e550 by goroutine 8:
  ase/etcd-20656.(*RecordingClient).AppendRespUnsafe()
      /work/client.go:18 +0x1a4
  ase/etcd-20656.TestRace_20656.func1()
      /work/verified_test.go:24 +0xad
  ase/etcd-20656.TestRace_20656.gowrap1()
      /work/verified_test.go:28 +0x41

Goroutine 14 (running) created at:
  ase/etcd-20656.TestRace_20656()
      /work/verified_test.go:22 +0x177
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/etcd-20656.TestRace_20656()
      /work/verified_test.go:22 +0x177
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-20656-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-20656-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-20656-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-20656-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-20656-bug .
# (then run as above, no --ssh flag)
```
