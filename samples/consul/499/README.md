# consul-499

| Field | Value |
|---|---|
| Project | consul |
| Reference | https://github.com/hashicorp/consul/pull/499 |
| Bug commit | `f126bb7381ae` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `consul/pool.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00011c010 by goroutine 69:
  ase/consul-499.(*ConnPool).reap()
      /work/verified_test.go:48 +0x52
  ase/consul-499.TestRace.func1()
      /work/verified_test.go:72 +0x6c

Previous write at 0x00c00011c010 by goroutine 70:
  ase/consul-499.(*ConnPool).Shutdown()
      /work/verified_test.go:41 +0x87
  ase/consul-499.TestRace.func2()
      /work/verified_test.go:79 +0x6c

Goroutine 69 (running) created at:
  ase/consul-499.TestRace()
      /work/verified_test.go:70 +0x21d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1203 +0x202

Goroutine 70 (finished) created at:
  ase/consul-499.TestRace()
      /work/verified_test.go:77 +0x7a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1203 +0x202
==================
--- FAIL: TestRace (1.45s)
    testing.go:1102: race detected during execution of test
FAIL
FAIL	ase/consul-499	1.463s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-consul-499-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-consul-499-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-consul-499-fix .
docker run --rm --memory=2g --cpus=1 gonb-consul-499-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-consul-499-bug .
# (then run as above, no --ssh flag)
```
