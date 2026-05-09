# etcd-1316

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/1316 |
| Bug commit | `da2ee9a90c46` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `etcdserver/stats/leader.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00018c008 by goroutine 15:
  ase/etcd-1316.(*FollowerStats).Succ()
      /work/leader.go:32 +0x1c9
  ase/etcd-1316.TestRace_1316.func1()
      /work/verified_test.go:21 +0xcd

Previous read at 0x00c00018c008 by goroutine 10:
  ase/etcd-1316.(*FollowerStats).Succ()
      /work/leader.go:30 +0x11d
  ase/etcd-1316.TestRace_1316.func1()
      /work/verified_test.go:21 +0xcd

Goroutine 15 (running) created at:
  ase/etcd-1316.TestRace_1316()
      /work/verified_test.go:18 +0x9d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 10 (running) created at:
  ase/etcd-1316.TestRace_1316()
      /work/verified_test.go:18 +0x9d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c00018c038 by goroutine 10:
  ase/etcd-1316.(*FollowerStats).Succ()
      /work/leader.go:31 +0x155
  ase/etcd-1316.TestRace_1316.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-1316-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-1316-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-1316-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-1316-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-1316-bug .
# (then run as above, no --ssh flag)
```
