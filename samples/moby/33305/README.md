# moby-33305

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/33305 |
| Bug commit | `d192db0d9350` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `daemon/logger/loggerutils/rotatefilewriter.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d4578 by goroutine 8:
  ase/moby-33305.(*RotateFileWriter).Write()
      /work/rotatefilewriter.go:36 +0x2a4
  ase/moby-33305.TestRace_33305.func1()
      /work/verified_test.go:36 +0x16c
  ase/moby-33305.TestRace_33305.gowrap3()
      /work/verified_test.go:38 +0x41

Previous read at 0x00c0000d4578 by goroutine 47:
  ase/moby-33305.(*RotateFileWriter).LogPath()
      /work/rotatefilewriter.go:47 +0xa7
  ase/moby-33305.TestRace_33305.func2()
      /work/verified_test.go:45 +0x9f

Goroutine 8 (running) created at:
  ase/moby-33305.TestRace_33305()
      /work/verified_test.go:33 +0x408
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 47 (finished) created at:
  ase/moby-33305.TestRace_33305()
      /work/verified_test.go:42 +0x564
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c000212000 by goroutine 32:
  os.(*File).Name()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-moby-33305-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-moby-33305-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-moby-33305-fix .
docker run --rm --memory=2g --cpus=1 gonb-moby-33305-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-moby-33305-bug .
# (then run as above, no --ssh flag)
```
