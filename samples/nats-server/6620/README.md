# nats-server-6620

| Field | Value |
|---|---|
| Project | nats-server |
| Reference | https://github.com/nats-io/nats-server/pull/6620 |
| Bug commit | `fcaef6fc3957` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `server/server.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0001e6018 by goroutine 10:
  ase/nats-server-6620.(*Server).rotateSysAccount()
      /work/server.go:23 +0xe6
  ase/nats-server-6620.TestRaceConfigureAccountsSysAccount.func1()
      /work/verified_test.go:28 +0x174
  ase/nats-server-6620.TestRaceConfigureAccountsSysAccount.gowrap1()
      /work/verified_test.go:30 +0x41

Previous read at 0x00c0001e6018 by goroutine 23:
  ase/nats-server-6620.(*Server).configureAccounts()
      /work/server.go:33 +0x50
  ase/nats-server-6620.TestRaceConfigureAccountsSysAccount.func2()
      /work/verified_test.go:37 +0x99

Goroutine 10 (running) created at:
  ase/nats-server-6620.TestRaceConfigureAccountsSysAccount()
      /work/verified_test.go:25 +0x186
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 23 (running) created at:
  ase/nats-server-6620.TestRaceConfigureAccountsSysAccount()
      /work/verified_test.go:34 +0x2d1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c000294008 by goroutine 18:
  ase/nats-server-6620.(*Server).configureAccounts()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-nats-server-6620-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-nats-server-6620-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-nats-server-6620-fix .
docker run --rm --memory=2g --cpus=1 gonb-nats-server-6620-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-nats-server-6620-bug .
# (then run as above, no --ssh flag)
```
