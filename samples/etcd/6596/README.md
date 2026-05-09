# etcd-6596

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/6596 |
| Bug commit | `b8079b7fc05f` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `lease/lessor.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000014268 by goroutine 15:
  ase/etcd-6596.TestRace_PR6596_LeaseTTL.func2()
      /work/verified_test.go:38 +0xa4

Previous read at 0x00c000014268 by goroutine 8:
  ase/etcd-6596.(*lessor).Renew()
      /work/lessor.go:55 +0x124
  ase/etcd-6596.TestRace_PR6596_LeaseTTL.func1()
      /work/verified_test.go:31 +0xbc

Goroutine 15 (running) created at:
  ase/etcd-6596.TestRace_PR6596_LeaseTTL()
      /work/verified_test.go:34 +0x3bb
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/etcd-6596.TestRace_PR6596_LeaseTTL()
      /work/verified_test.go:27 +0x4ea
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c000014268 by goroutine 23:
  ase/etcd-6596.TestRace_PR6596_LeaseTTL.func2()
      /work/verified_test.go:38 +0xa4

Previous write at 0x00c000014268 by goroutine 9:
  ase/etcd-6596.TestRace_PR6596_LeaseTTL.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-6596-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-6596-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-6596-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-6596-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-6596-bug .
# (then run as above, no --ssh flag)
```
