# consul-1214

| Field | Value |
|---|---|
| Project | consul |
| Reference | https://github.com/hashicorp/consul/pull/1214 |
| Bug commit | `f41b79eff2af` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `command/lock.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0001904a8 by goroutine 9:
  ase/consul-1214.(*LockCommand).killChild()
      /work/verified_test.go:41 +0x73
  ase/consul-1214.TestRace.func2()
      /work/verified_test.go:70 +0x66

Previous write at 0x00c0001904a8 by goroutine 8:
  ase/consul-1214.(*LockCommand).startChild()
      /work/verified_test.go:32 +0x84
  ase/consul-1214.TestRace.func1()
      /work/verified_test.go:62 +0x6c

Goroutine 9 (running) created at:
  ase/consul-1214.TestRace()
      /work/verified_test.go:68 +0x7a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1203 +0x202

Goroutine 8 (finished) created at:
  ase/consul-1214.TestRace()
      /work/verified_test.go:60 +0x184
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1203 +0x202
==================
==================
WARNING: DATA RACE
Read at 0x00c000116000 by goroutine 10:
  ase/consul-1214.(*LockCommand).killChild()
      /work/verified_test.go:44 +0xa4
  ase/consul-1214.TestRace.func2()
      /work/verified_test.go:70 +0x66

Previous write at 0x00c000116000 by goroutine 8:
  ase/consul-1214.(*LockCommand).startChild()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-consul-1214-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-consul-1214-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-consul-1214-fix .
docker run --rm --memory=2g --cpus=1 gonb-consul-1214-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-consul-1214-bug .
# (then run as above, no --ssh flag)
```
