# consul-3162

| Field | Value |
|---|---|
| Project | consul |
| Reference | https://github.com/hashicorp/consul/pull/3162 |
| Bug commit | `5473255f9816` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `agent/agent.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000146038 by goroutine 68:
  ase/consul-3162.TestShutdownRace.func2()
      /work/verified_test.go:52 +0x75

Previous write at 0x00c000146038 by goroutine 8:
  ase/consul-3162.TestShutdownRace.func1()
      /work/verified_test.go:41 +0xc7

Goroutine 68 (running) created at:
  ase/consul-3162.TestShutdownRace()
      /work/verified_test.go:49 +0xca
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1203 +0x202

Goroutine 8 (finished) created at:
  ase/consul-3162.TestShutdownRace()
      /work/verified_test.go:34 +0x78
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1203 +0x202
==================
--- FAIL: TestShutdownRace (0.97s)
    testing.go:1102: race detected during execution of test
FAIL
FAIL	ase/consul-3162	0.987s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-consul-3162-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-consul-3162-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-consul-3162-fix .
docker run --rm --memory=2g --cpus=1 gonb-consul-3162-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-consul-3162-bug .
# (then run as above, no --ssh flag)
```
