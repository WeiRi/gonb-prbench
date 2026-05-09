# dns-780

| Field | Value |
|---|---|
| Project | dns |
| Reference | https://github.com/miekg/dns/pull/780 |
| Bug commit | `c0283a202831` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `server.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00012c8d8 by goroutine 9:
  ase/dns-780.(*Server).ShutdownContext()
      /work/server.go:40 +0x98
  ase/dns-780.TestRace_dns780_started_vs_shutdown.func2()
      /work/verified_test.go:24 +0x99

Previous read at 0x00c00012c8d8 by goroutine 8:
  ase/dns-780.(*Server).ReadTCP()
      /work/server.go:27 +0x37
  ase/dns-780.TestRace_dns780_started_vs_shutdown.func1()
      /work/verified_test.go:18 +0x99

Goroutine 9 (running) created at:
  ase/dns-780.TestRace_dns780_started_vs_shutdown()
      /work/verified_test.go:21 +0x20c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/dns-780.TestRace_dns780_started_vs_shutdown()
      /work/verified_test.go:15 +0x164
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c00012c8e0 by goroutine 9:
  ase/dns-780.(*Server).ShutdownContext()
      /work/server.go:41 +0xd2
  ase/dns-780.TestRace_dns780_started_vs_shutdown.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-dns-780-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-dns-780-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-dns-780-fix .
docker run --rm --memory=2g --cpus=1 gonb-dns-780-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-dns-780-bug .
# (then run as above, no --ssh flag)
```
