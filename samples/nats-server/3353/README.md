# nats-server-3353

| Field | Value |
|---|---|
| Project | nats-server |
| Reference | https://github.com/nats-io/nats-server/pull/3353 |
| Bug commit | `9a92d10cc906` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `server/filestore.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000012140 by goroutine 9:
  ase/nats-server-3353.(*msgBlock).writePerSubject()
      /work/filestore.go:20 +0x73
  ase/nats-server-3353.TestRacePopulateGlobalPerSubjectInfo.func1()
      /work/verified_test.go:32 +0x179
  ase/nats-server-3353.TestRacePopulateGlobalPerSubjectInfo.gowrap1()
      /work/verified_test.go:34 +0x41

Previous read at 0x00c000012140 by goroutine 17:
  ase/nats-server-3353.(*msgBlock).readPerSubjectInfo()
      /work/filestore.go:28 +0x16a
  ase/nats-server-3353.(*fileStore).populateGlobalPerSubjectInfo()
      /work/filestore.go:36 +0x162
  ase/nats-server-3353.TestRacePopulateGlobalPerSubjectInfo.func2()
      /work/verified_test.go:41 +0xea

Goroutine 9 (running) created at:
  ase/nats-server-3353.TestRacePopulateGlobalPerSubjectInfo()
      /work/verified_test.go:29 +0xf7
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 17 (finished) created at:
  ase/nats-server-3353.TestRacePopulateGlobalPerSubjectInfo()
      /work/verified_test.go:38 +0x235
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-nats-server-3353-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-nats-server-3353-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-nats-server-3353-fix .
docker run --rm --memory=2g --cpus=1 gonb-nats-server-3353-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-nats-server-3353-bug .
# (then run as above, no --ssh flag)
```
