# consul-1158

| Field | Value |
|---|---|
| Project | consul |
| Reference | https://github.com/hashicorp/consul/pull/1158 |
| Bug commit | `3846dea3706f` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `command/lock.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0001c0040 by goroutine 11:
  ase/consul-1158.(*LockCommand).startChild()
      /work/verified_test.go:29 +0x44
  ase/consul-1158.TestRace.func1()
      /work/verified_test.go:64 +0x6c

Previous write at 0x00c0001c0040 by goroutine 8:
  ase/consul-1158.(*LockCommand).startChild()
      /work/verified_test.go:29 +0x44
  ase/consul-1158.TestRace.func1()
      /work/verified_test.go:64 +0x6c

Goroutine 11 (running) created at:
  ase/consul-1158.TestRace()
      /work/verified_test.go:62 +0x78
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1203 +0x202

Goroutine 8 (finished) created at:
  ase/consul-1158.TestRace()
      /work/verified_test.go:62 +0x78
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1203 +0x202
==================
--- FAIL: TestRace (1.40s)
    testing.go:1102: race detected during execution of test
FAIL
FAIL	ase/consul-1158	1.420s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-consul-1158-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-consul-1158-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-consul-1158-fix .
docker run --rm --memory=2g --cpus=1 gonb-consul-1158-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-consul-1158-bug .
# (then run as above, no --ssh flag)
```
