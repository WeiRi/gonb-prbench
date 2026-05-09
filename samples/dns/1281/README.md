# dns-1281

| Field | Value |
|---|---|
| Project | dns |
| Reference | https://github.com/miekg/dns/pull/1281 |
| Bug commit | `af0c865ab359` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `client.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00009a080 by goroutine 9:
  ase/dns-1281.(*Client).Dial()
      /work/client.go:22 +0xdc
  ase/dns-1281.TestRace_dns1281_Client_Dialer.func2()
      /work/verified_test.go:25 +0xed

Previous write at 0x00c00009a080 by goroutine 8:
  ase/dns-1281.(*Client).ExchangeContext()
      /work/client.go:33 +0x13c
  ase/dns-1281.TestRace_dns1281_Client_Dialer.func1()
      /work/verified_test.go:19 +0x11d

Goroutine 9 (running) created at:
  ase/dns-1281.TestRace_dns1281_Client_Dialer()
      /work/verified_test.go:22 +0x1e4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/dns-1281.TestRace_dns1281_Client_Dialer()
      /work/verified_test.go:16 +0x13c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c00009a080 by goroutine 9:
  ase/dns-1281.(*Client).Dial()
      /work/client.go:25 +0xec
  ase/dns-1281.TestRace_dns1281_Client_Dialer.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-dns-1281-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-dns-1281-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-dns-1281-fix .
docker run --rm --memory=2g --cpus=1 gonb-dns-1281-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-dns-1281-bug .
# (then run as above, no --ssh flag)
```
