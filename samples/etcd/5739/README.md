# etcd-5739

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/5739 |
| Bug commit | `27ef4baa9c35` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `etcdserver/apply_auth.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000122128 by goroutine 11:
  ase/etcd-5739.(*authApplierV3).Apply()
      /work/apply_auth.go:16 +0xd2
  ase/etcd-5739.TestRace_PR5739_AuthApplierUserUnlocked.func1()
      /work/verified_test.go:20 +0xca
  ase/etcd-5739.TestRace_PR5739_AuthApplierUserUnlocked.gowrap1()
      /work/verified_test.go:22 +0x41

Previous write at 0x00c000122128 by goroutine 15:
  ase/etcd-5739.(*authApplierV3).Apply()
      /work/apply_auth.go:16 +0xd2
  ase/etcd-5739.TestRace_PR5739_AuthApplierUserUnlocked.func1()
      /work/verified_test.go:20 +0xca
  ase/etcd-5739.TestRace_PR5739_AuthApplierUserUnlocked.gowrap1()
      /work/verified_test.go:22 +0x41

Goroutine 11 (running) created at:
  ase/etcd-5739.TestRace_PR5739_AuthApplierUserUnlocked()
      /work/verified_test.go:17 +0xcb
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 15 (finished) created at:
  ase/etcd-5739.TestRace_PR5739_AuthApplierUserUnlocked()
      /work/verified_test.go:17 +0xcb
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-5739-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-5739-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-5739-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-5739-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-5739-bug .
# (then run as above, no --ssh flag)
```
