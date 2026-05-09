# etcd-4876

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/4876 |
| Bug commit | `1637b371320a` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `clientv3/client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0001182e8 by goroutine 23:
  ase/etcd-4876.(*kV).switchRemote()
      /work/kv.go:27 +0x3c
  ase/etcd-4876.TestRace_PR4876_SwitchRemoteUnlocked.func2()
      /work/verified_test.go:27 +0x99

Previous write at 0x00c0001182e8 by goroutine 11:
  ase/etcd-4876.(*kV).switchRemote()
      /work/kv.go:30 +0xa4
  ase/etcd-4876.TestRace_PR4876_SwitchRemoteUnlocked.func2()
      /work/verified_test.go:27 +0x99

Goroutine 23 (running) created at:
  ase/etcd-4876.TestRace_PR4876_SwitchRemoteUnlocked()
      /work/verified_test.go:24 +0x111
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 11 (running) created at:
  ase/etcd-4876.TestRace_PR4876_SwitchRemoteUnlocked()
      /work/verified_test.go:24 +0x111
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR4876_SwitchRemoteUnlocked (0.19s)
FAIL
FAIL	ase/etcd-4876	0.201s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-4876-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-4876-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-4876-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-4876-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-4876-bug .
# (then run as above, no --ssh flag)
```
