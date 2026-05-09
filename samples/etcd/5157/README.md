# etcd-5157

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/5157 |
| Bug commit | `69bc0f76bc04` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `etcdserver/stats/server.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000161f0 by goroutine 9:
  ase/etcd-5157.(*ServerStats).GetState()
      /work/server.go:36 +0xc4
  ase/etcd-5157.TestRace_5157.func2()
      /work/verified_test.go:33 +0xb8

Previous write at 0x00c0000161f0 by goroutine 8:
  ase/etcd-5157.(*ServerStats).SetState()
      /work/server.go:29 +0xda
  ase/etcd-5157.TestRace_5157.func1()
      /work/verified_test.go:25 +0xd2

Goroutine 9 (running) created at:
  ase/etcd-5157.TestRace_5157()
      /work/verified_test.go:30 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/etcd-5157.TestRace_5157()
      /work/verified_test.go:22 +0x244
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c0000161f8 by goroutine 9:
  ase/etcd-5157.(*ServerStats).GetState()
      /work/server.go:36 +0xce
  ase/etcd-5157.TestRace_5157.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-5157-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-5157-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-5157-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-5157-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-5157-bug .
# (then run as above, no --ssh flag)
```
