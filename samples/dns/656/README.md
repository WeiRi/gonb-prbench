# dns-656

| Field | Value |
|---|---|
| Project | dns |
| Reference | https://github.com/miekg/dns/pull/656 |
| Bug commit | `5b169d1842fb` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000016348 by goroutine 29:
  github.com/miekg/dns.(*Conn).WriteMsg()
      /work/client.go:364 +0x2b9
  github.com/miekg/dns.TestRaceConnSharedMutableState.func2()
      /work/race_656_capture_test.go:45 +0x7c
  github.com/miekg/dns.TestRaceConnSharedMutableState.gowrap4()
      /work/race_656_capture_test.go:54 +0x41

Previous write at 0x00c000016348 by goroutine 14:
  github.com/miekg/dns.(*Conn).WriteMsg()
      /work/client.go:364 +0x2b9
  github.com/miekg/dns.TestRaceConnSharedMutableState.func2()
      /work/race_656_capture_test.go:45 +0x7c
  github.com/miekg/dns.TestRaceConnSharedMutableState.gowrap4()
      /work/race_656_capture_test.go:54 +0x41

Goroutine 29 (running) created at:
  github.com/miekg/dns.TestRaceConnSharedMutableState()
      /work/race_656_capture_test.go:40 +0x430
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44

Goroutine 14 (running) created at:
  github.com/miekg/dns.TestRaceConnSharedMutableState()
      /work/race_656_capture_test.go:40 +0x430
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
==================
==================
WARNING: DATA RACE
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-dns-656-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-dns-656-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-dns-656-fix .
docker run --rm --memory=2g --cpus=1 gonb-dns-656-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-dns-656-bug .
# (then run as above, no --ssh flag)
```
