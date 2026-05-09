# etcd-6704

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/6704 |
| Bug commit | `21e65eec0890` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `proxy/grpcproxy/watcher_group.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000014258 by goroutine 8:
  ase/etcd-6704.(*watcherGroup).broadcast()
      /work/watcher_group.go:29 +0x44
  ase/etcd-6704.TestRace_PR6704_GroupRevUnlocked.func1()
      /work/verified_test.go:25 +0xa9
  ase/etcd-6704.TestRace_PR6704_GroupRevUnlocked.gowrap1()
      /work/verified_test.go:27 +0x41

Previous read at 0x00c000014258 by goroutine 13:
  ase/etcd-6704.(*watcherGroups).maybeJoinWatcherSingle()
      /work/watcher_group.go:48 +0x12f
  ase/etcd-6704.TestRace_PR6704_GroupRevUnlocked.func2()
      /work/verified_test.go:31 +0x130
  ase/etcd-6704.TestRace_PR6704_GroupRevUnlocked.gowrap2()
      /work/verified_test.go:36 +0x41

Goroutine 8 (running) created at:
  ase/etcd-6704.TestRace_PR6704_GroupRevUnlocked()
      /work/verified_test.go:22 +0x2f1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 13 (finished) created at:
  ase/etcd-6704.TestRace_PR6704_GroupRevUnlocked()
      /work/verified_test.go:28 +0x1ab
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR6704_GroupRevUnlocked (0.02s)
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-6704-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-6704-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-6704-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-6704-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-6704-bug .
# (then run as above, no --ssh flag)
```
