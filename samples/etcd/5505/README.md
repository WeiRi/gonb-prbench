# etcd-5505

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/5505 |
| Bug commit | `310ebdd3e15c` |
| Category | channel_misuse |
| Oracle | RACE |
| Primary diff file | `etcdserver/api/v3rpc/watch.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000ce010 by goroutine 41:
  runtime.chansend()
      /usr/local/go/src/runtime/chan.go:160 +0x0
  ase/etcd-5505.(*serverWatchStream).recvLoop()
      /work/watch.go:26 +0x131
  ase/etcd-5505.TestRace_PR5505_CtrlChanClose.func1.1()
      /work/verified_test.go:25 +0x48

Previous write at 0x00c0000ce010 by goroutine 8:
  runtime.closechan()
      /usr/local/go/src/runtime/chan.go:357 +0x0
  ase/etcd-5505.(*serverWatchStream).close.func1()
      /work/watch.go:48 +0x45
  sync.(*Once).doSlow()
      /usr/local/go/src/sync/once.go:74 +0xf0
  sync.(*Once).Do()
      /usr/local/go/src/sync/once.go:65 +0x44
  ase/etcd-5505.(*serverWatchStream).close()
      /work/watch.go:47 +0x64
  ase/etcd-5505.TestRace_PR5505_CtrlChanClose.func1.2()
      /work/verified_test.go:31 +0x44
  ase/etcd-5505.TestRace_PR5505_CtrlChanClose.func1()
      /work/verified_test.go:32 +0x264

Goroutine 41 (running) created at:
  ase/etcd-5505.TestRace_PR5505_CtrlChanClose.func1()
      /work/verified_test.go:23 +0x1e4

Goroutine 8 (running) created at:
  ase/etcd-5505.TestRace_PR5505_CtrlChanClose()
      /work/verified_test.go:19 +0x64
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-5505-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-5505-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-5505-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-5505-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-5505-bug .
# (then run as above, no --ssh flag)
```
