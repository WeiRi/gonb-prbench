# go-ethereum-20204

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/20204 |
| Bug commit | `b0b277525cb4` |
| Category | anonymous_function |
| Oracle | RACE |
| Primary diff file | `eth/downloader/downloader.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00009a078 by goroutine 8:
  ase/go-ethereum-20204.(*Downloader).processFastSyncContent.func1()
      /work/downloader.go:32 +0x2e

Previous write at 0x00c00009a078 by goroutine 7:
  ase/go-ethereum-20204.(*Downloader).processFastSyncContent()
      /work/downloader.go:38 +0x1d2
  ase/go-ethereum-20204.TestRace_20204_downloader_capture()
      /work/verified_test.go:12 +0x26
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/go-ethereum-20204.(*Downloader).processFastSyncContent()
      /work/downloader.go:31 +0x13c
  ase/go-ethereum-20204.TestRace_20204_downloader_capture()
      /work/verified_test.go:12 +0x26
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 7 (running) created at:
  testing.(*T).Run()
      /usr/local/go/src/testing/testing.go:1742 +0x825
  testing.runTests.func1()
      /usr/local/go/src/testing/testing.go:2161 +0x85
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.runTests()
      /usr/local/go/src/testing/testing.go:2159 +0x8be
  testing.(*M).Run()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-go-ethereum-20204-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-20204-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-go-ethereum-20204-fix .
docker run --rm --memory=2g --cpus=1 gonb-go-ethereum-20204-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-go-ethereum-20204-bug .
# (then run as above, no --ssh flag)
```
