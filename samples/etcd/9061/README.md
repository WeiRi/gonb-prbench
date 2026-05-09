# etcd-9061

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/9061 |
| Bug commit | `3dd1c1b53c65` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `clientv3/leasing/kv.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d4550 by goroutine 9:
  ase/etcd-9061.TestRace_PR9061_WaitSession.func2()
      /work/verified_test.go:32 +0xc9

Previous write at 0x00c0000d4550 by goroutine 11:
  ase/etcd-9061.TestRace_PR9061_WaitSession.func2()
      /work/verified_test.go:32 +0xc9

Goroutine 9 (running) created at:
  ase/etcd-9061.TestRace_PR9061_WaitSession()
      /work/verified_test.go:29 +0xfd
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 11 (finished) created at:
  ase/etcd-9061.TestRace_PR9061_WaitSession()
      /work/verified_test.go:29 +0xfd
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0000d4550 by goroutine 9:
  ase/etcd-9061.TestRace_PR9061_WaitSession.func2()
      /work/verified_test.go:32 +0xc9

Previous read at 0x00c0000d4550 by goroutine 14:
  ase/etcd-9061.(*leasingKV).waitSession()
      /work/kv.go:20 +0x10d
  ase/etcd-9061.TestRace_PR9061_WaitSession.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-9061-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-9061-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-9061-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-9061-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-9061-bug .
# (then run as above, no --ssh flag)
```
