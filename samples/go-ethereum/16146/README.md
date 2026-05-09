# go-ethereum-16146

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/16146 |
| Bug commit | `5603715c0699` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `whisper/whisperv6/peer.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d4550 by goroutine 8:
  ase/go-ethereum-16146.(*Peer).setBloomFilter()
      /work/peer.go:58 +0x10f
  ase/go-ethereum-16146.TestRace_16146_bloomFilter_concurrent.func1()
      /work/verified_test.go:26 +0x107

Previous read at 0x00c0000d4550 by goroutine 9:
  ase/go-ethereum-16146.(*Peer).bloomMatch()
      /work/peer.go:52 +0x133
  ase/go-ethereum-16146.TestRace_16146_bloomFilter_concurrent.func2()
      /work/verified_test.go:32 +0xe2

Goroutine 8 (running) created at:
  ase/go-ethereum-16146.TestRace_16146_bloomFilter_concurrent()
      /work/verified_test.go:19 +0x1b0
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (finished) created at:
  ase/go-ethereum-16146.TestRace_16146_bloomFilter_concurrent()
      /work/verified_test.go:29 +0x29c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0000d4568 by goroutine 8:
  ase/go-ethereum-16146.(*Peer).setBloomFilter()
      /work/peer.go:59 +0x184
  ase/go-ethereum-16146.TestRace_16146_bloomFilter_concurrent.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-go-ethereum-16146-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-16146-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-go-ethereum-16146-fix .
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-16146-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-go-ethereum-16146-bug .
# (then run as above, no --ssh flag)
```
