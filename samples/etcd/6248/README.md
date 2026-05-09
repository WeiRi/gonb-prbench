# etcd-6248

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/6248 |
| Bug commit | `028b9540520d` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `mvcc/kvstore.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c420016110 by goroutine 11:
  _/work.TestRace_6248.func1()
      /work/verified_test.go:7 +0x6b

Previous write at 0x00c420016110 by goroutine 8:
  _/work.TestRace_6248.func1()
      /work/verified_test.go:7 +0x6b

Goroutine 11 (running) created at:
  _/work.TestRace_6248()
      /work/verified_test.go:13 +0xdb
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:777 +0x16d

Goroutine 8 (finished) created at:
  _/work.TestRace_6248()
      /work/verified_test.go:13 +0xdb
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:777 +0x16d
==================
==================
WARNING: DATA RACE
Write at 0x00c420016110 by goroutine 10:
  _/work.TestRace_6248.func1()
      /work/verified_test.go:7 +0x6b

Previous write at 0x00c420016110 by goroutine 8:
  _/work.TestRace_6248.func1()
      /work/verified_test.go:7 +0x6b

Goroutine 10 (running) created at:
  _/work.TestRace_6248()
      /work/verified_test.go:13 +0xdb
  testing.tRunner()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-6248-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-6248-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-6248-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-6248-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-6248-bug .
# (then run as above, no --ssh flag)
```
