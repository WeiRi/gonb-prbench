# nomad-14538

| Field | Value |
|---|---|
| Project | nomad |
| Reference | https://github.com/hashicorp/nomad/pull/14538 |
| Bug commit | `39a3fd652c01` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `client/logmon/logging/rotator.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000016228 by goroutine 12:
  ase/nomad-14538.(*FileRotator).purgeOldFiles()
      /work/rotator.go:26 +0xc9
  ase/nomad-14538.TestRace_OldestLogFileIdx.func1()
      /work/verified_test.go:25 +0xa4
  ase/nomad-14538.TestRace_OldestLogFileIdx.gowrap2()
      /work/verified_test.go:27 +0x41

Previous write at 0x00c000016228 by goroutine 8:
  ase/nomad-14538.(*FileRotator).purgeOldFiles()
      /work/rotator.go:26 +0xc9
  ase/nomad-14538.TestRace_OldestLogFileIdx.func1()
      /work/verified_test.go:25 +0xa4
  ase/nomad-14538.TestRace_OldestLogFileIdx.gowrap2()
      /work/verified_test.go:27 +0x41

Goroutine 12 (running) created at:
  ase/nomad-14538.TestRace_OldestLogFileIdx()
      /work/verified_test.go:22 +0x257
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/nomad-14538.TestRace_OldestLogFileIdx()
      /work/verified_test.go:22 +0x257
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-nomad-14538-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-nomad-14538-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-nomad-14538-fix .
docker run --rm --memory=2g --cpus=1 gonb-nomad-14538-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-nomad-14538-bug .
# (then run as above, no --ssh flag)
```
