# etcd-3077

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/3077 |
| Bug commit | `235aef53650f` |
| Category | order_violation |
| Oracle | RACE |
| Primary diff file | `etcdserver/raft.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00011a000 by goroutine 28:
  ase/etcd-3077.(*EtcdServer).run.func1()
      /work/raft.go:24 +0x5a

Previous read at 0x00c00011a000 by goroutine 14:
  ase/etcd-3077.(*EtcdServer).run()
      /work/raft.go:30 +0xa8
  ase/etcd-3077.TestRace_PR3077_RaftInitOrdering.func1()
      /work/verified_test.go:19 +0x10a

Goroutine 28 (running) created at:
  ase/etcd-3077.(*EtcdServer).run()
      /work/raft.go:22 +0x8d
  ase/etcd-3077.TestRace_PR3077_RaftInitOrdering.func1()
      /work/verified_test.go:19 +0x10a

Goroutine 14 (finished) created at:
  ase/etcd-3077.TestRace_PR3077_RaftInitOrdering()
      /work/verified_test.go:16 +0x64
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c00011a008 by goroutine 28:
  ase/etcd-3077.(*EtcdServer).run.func1()
      /work/raft.go:25 +0xc4

Previous read at 0x00c00011a008 by goroutine 14:
  ase/etcd-3077.(*EtcdServer).run()
      /work/raft.go:31 +0xc5
  ase/etcd-3077.TestRace_PR3077_RaftInitOrdering.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-3077-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-3077-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-3077-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-3077-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-3077-bug .
# (then run as above, no --ssh flag)
```
