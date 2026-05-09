# etcd-2968

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/2968 |
| Bug commit | `684c7213076e` |
| Category | order_violation |
| Oracle | RACE |
| Primary diff file | `rafthttp/peer.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000012020 by goroutine 18:
  ase/etcd-2968.startPeer.func1()
      /work/peer.go:20 +0x59

Previous read at 0x00c000012020 by goroutine 8:
  ase/etcd-2968.(*peer).stop()
      /work/peer.go:29 +0xa6
  ase/etcd-2968.TestRace_PR2968_MsgAppReaderInit.func1()
      /work/verified_test.go:19 +0x86

Goroutine 18 (running) created at:
  ase/etcd-2968.startPeer()
      /work/peer.go:18 +0x124
  ase/etcd-2968.TestRace_PR2968_MsgAppReaderInit.func1()
      /work/verified_test.go:18 +0x71

Goroutine 8 (finished) created at:
  ase/etcd-2968.TestRace_PR2968_MsgAppReaderInit()
      /work/verified_test.go:16 +0x64
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR2968_MsgAppReaderInit (0.00s)
FAIL
FAIL	ase/etcd-2968	0.019s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-2968-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-2968-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-2968-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-2968-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-2968-bug .
# (then run as above, no --ssh flag)
```
