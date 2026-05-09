# rclone-140

| Field | Value |
|---|---|
| Project | rclone |
| Reference | https://github.com/rclone/rclone/pull/140 |
| Bug commit | `34193fd8d98b` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pacer/pacer.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000d6530 by goroutine 10:
  ase/rclone-140.(*Pacer).endCall()
      /work/pacer.go:49 +0xc8
  ase/rclone-140.TestRace_140_PacerFields.func3()
      /work/verified_test.go:29 +0xc9

Previous write at 0x00c0000d6530 by goroutine 8:
  ase/rclone-140.(*Pacer).SetMinSleep()
      /work/pacer.go:24 +0xe7
  ase/rclone-140.TestRace_140_PacerFields.func1()
      /work/verified_test.go:17 +0x99

Goroutine 10 (running) created at:
  ase/rclone-140.TestRace_140_PacerFields()
      /work/verified_test.go:26 +0x35c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/rclone-140.TestRace_140_PacerFields()
      /work/verified_test.go:14 +0x204
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0000d6538 by goroutine 9:
  ase/rclone-140.(*Pacer).SetRetries()
      /work/pacer.go:43 +0xac
  ase/rclone-140.TestRace_140_PacerFields.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-rclone-140-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-rclone-140-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-rclone-140-fix .
docker run --rm --memory=2g --cpus=1 gonb-rclone-140-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-rclone-140-bug .
# (then run as above, no --ssh flag)
```
