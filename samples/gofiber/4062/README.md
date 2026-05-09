# gofiber-4062

| Field | Value |
|---|---|
| Project | gofiber |
| Reference | https://github.com/gofiber/fiber/pull/4062 |
| Bug commit | `b45321db1cd0` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `ctx.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000096080 by goroutine 9:
  ase/gofiber-4062.(*DefaultCtx).release()
      /work/ctx.go:37 +0x84
  ase/gofiber-4062.TestRace_PR4062_ValueAfterRelease.func2()
      /work/verified_test.go:32 +0x7d

Previous read at 0x00c000096080 by goroutine 8:
  ase/gofiber-4062.(*DefaultCtx).Value()
      /work/ctx.go:32 +0xd0
  ase/gofiber-4062.TestRace_PR4062_ValueAfterRelease.func1()
      /work/verified_test.go:23 +0xcb

Goroutine 9 (running) created at:
  ase/gofiber-4062.TestRace_PR4062_ValueAfterRelease()
      /work/verified_test.go:30 +0x2a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/gofiber-4062.TestRace_PR4062_ValueAfterRelease()
      /work/verified_test.go:20 +0x2e4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_PR4062_ValueAfterRelease (0.01s)
FAIL
FAIL	ase/gofiber-4062	0.029s
FAIL
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-gofiber-4062-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-gofiber-4062-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-gofiber-4062-fix .
docker run --rm --memory=2g --cpus=1 gonb-gofiber-4062-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-gofiber-4062-bug .
# (then run as above, no --ssh flag)
```
