# dns-459

| Field | Value |
|---|---|
| Project | dns |
| Reference | https://github.com/miekg/dns/pull/459 |
| Bug commit | `74ec3b243352` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00029a020 by goroutine 9:
  github.com/miekg/dns.TestRaceExchangeSharedResponseError.func2()
      /work/race_459_capture_test.go:40 +0xb3
  github.com/miekg/dns.TestRaceExchangeSharedResponseError.gowrap3()
      /work/race_459_capture_test.go:46 +0x41

Previous write at 0x00c00029a020 by goroutine 84:
  github.com/miekg/dns.TestRaceExchangeSharedResponseError.func2()
      /work/race_459_capture_test.go:40 +0xca
  github.com/miekg/dns.TestRaceExchangeSharedResponseError.gowrap3()
      /work/race_459_capture_test.go:46 +0x41

Goroutine 9 (running) created at:
  github.com/miekg/dns.TestRaceExchangeSharedResponseError()
      /work/race_459_capture_test.go:35 +0x284
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 84 (running) created at:
  github.com/miekg/dns.TestRaceExchangeSharedResponseError()
      /work/race_459_capture_test.go:35 +0x284
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c00029a020 by goroutine 9:
  github.com/miekg/dns.TestRaceExchangeSharedResponseError.func2()
      /work/race_459_capture_test.go:40 +0xca
  github.com/miekg/dns.TestRaceExchangeSharedResponseError.gowrap3()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-dns-459-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-dns-459-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-dns-459-fix .
docker run --rm --memory=2g --cpus=1 gonb-dns-459-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-dns-459-bug .
# (then run as above, no --ssh flag)
```
