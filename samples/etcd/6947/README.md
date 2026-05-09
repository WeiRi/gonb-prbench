# etcd-6947

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/6947 |
| Bug commit | `b9b14b15d624` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `proxy/grpcproxy/cache/store.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d4570 by goroutine 22:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/etcd-6947.(*Cache).Add()
      /work/store.go:27 +0x44b
  ase/etcd-6947.TestRace_PR6947_SizeUnlocked.func1()
      /work/verified_test.go:22 +0x17c
  ase/etcd-6947.TestRace_PR6947_SizeUnlocked.gowrap1()
      /work/verified_test.go:24 +0x41

Previous read at 0x00c0000d4570 by goroutine 13:
  ase/etcd-6947.(*Cache).Size()
      /work/store.go:32 +0xda
  ase/etcd-6947.TestRace_PR6947_SizeUnlocked.func2()
      /work/verified_test.go:28 +0xab

Goroutine 22 (running) created at:
  ase/etcd-6947.TestRace_PR6947_SizeUnlocked()
      /work/verified_test.go:19 +0x37c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 13 (finished) created at:
  ase/etcd-6947.TestRace_PR6947_SizeUnlocked()
      /work/verified_test.go:25 +0x231
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR6947_SizeUnlocked (0.06s)
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-6947-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-6947-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-6947-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-6947-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-6947-bug .
# (then run as above, no --ssh flag)
```
