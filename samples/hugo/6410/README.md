# hugo-6410

| Field | Value |
|---|---|
| Project | hugo |
| Reference | https://github.com/gohugoio/hugo/pull/6410 |
| Bug commit | `0d7b05be4cb2` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `helpers/general.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x000000717338 by goroutine 220:
  ase/hugo-6410.InitLoggers()
      /work/helpers_general.go:53 +0x14c
  ase/hugo-6410.TestRaceInitLoggers.func2()
      /work/verified_test.go:23 +0x71

Previous read at 0x000000717338 by goroutine 89:
  ase/hugo-6410.TestRaceInitLoggers.func1()
      /work/verified_test.go:16 +0x12f
  ase/hugo-6410.TestRaceInitLoggers.gowrap1()
      /work/verified_test.go:17 +0x41

Goroutine 220 (running) created at:
  ase/hugo-6410.TestRaceInitLoggers()
      /work/verified_test.go:21 +0x174
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 89 (running) created at:
  ase/hugo-6410.TestRaceInitLoggers()
      /work/verified_test.go:14 +0x70
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x000000717338 by goroutine 234:
  ase/hugo-6410.InitLoggers()
      /work/helpers_general.go:53 +0x14c
  ase/hugo-6410.TestRaceInitLoggers.func2()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-hugo-6410-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-hugo-6410-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-hugo-6410-fix .
docker run --rm --memory=2g --cpus=1 gonb-hugo-6410-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-hugo-6410-bug .
# (then run as above, no --ssh flag)
```
