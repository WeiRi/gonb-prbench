# gofiber-3985

| Field | Value |
|---|---|
| Project | gofiber |
| Reference | https://github.com/gofiber/fiber/pull/3985 |
| Bug commit | `80d999ef7e1c` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `middleware/cache/cache.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000be180 by goroutine 8:
  ase/gofiber-3985.(*item).updateData()
      /work/cache.go:17 +0xf3
  ase/gofiber-3985.TestRace_3985.func1()
      /work/verified_test.go:18 +0xc3
  ase/gofiber-3985.TestRace_3985.gowrap1()
      /work/verified_test.go:20 +0x41

Previous read at 0x00c0000be180 by goroutine 12:
  ase/gofiber-3985.(*item).readData()
      /work/cache.go:23 +0xa4
  ase/gofiber-3985.TestRace_3985.func2()
      /work/verified_test.go:24 +0x9b

Goroutine 8 (running) created at:
  ase/gofiber-3985.TestRace_3985()
      /work/verified_test.go:15 +0x1fc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 12 (finished) created at:
  ase/gofiber-3985.TestRace_3985()
      /work/verified_test.go:21 +0x2a4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x00c0000be180 by goroutine 8:
  ase/gofiber-3985.(*item).updateData()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-gofiber-3985-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-gofiber-3985-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-gofiber-3985-fix .
docker run --rm --memory=2g --cpus=1 gonb-gofiber-3985-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-gofiber-3985-bug .
# (then run as above, no --ssh flag)
```
