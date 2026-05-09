# etcd-20192

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/20192 |
| Bug commit | `88feab3eecbe` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `client/v3/client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00001c330 by goroutine 34:
  ase/etcd-20192.(*Client).SetLogger()
      /work/client.go:16 +0x124
  ase/etcd-20192.TestRace_20192.func1()
      /work/verified_test.go:30 +0x11f

Previous read at 0x00c00001c330 by goroutine 17:
  ase/etcd-20192.(*Client).GetLogger()
      /work/client.go:21 +0xd0
  ase/etcd-20192.TestRace_20192.func2()
      /work/verified_test.go:38 +0xc8

Goroutine 34 (running) created at:
  ase/etcd-20192.TestRace_20192()
      /work/verified_test.go:25 +0x31c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 17 (finished) created at:
  ase/etcd-20192.TestRace_20192()
      /work/verified_test.go:33 +0x1d1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c00001c330 by goroutine 30:
  ase/etcd-20192.(*Client).SetLogger()
      /work/client.go:16 +0x124
  ase/etcd-20192.TestRace_20192.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-20192-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-20192-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-20192-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-20192-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-20192-bug .
# (then run as above, no --ssh flag)
```
