# hugo-12413

| Field | Value |
|---|---|
| Project | hugo |
| Reference | https://github.com/gohugoio/hugo/pull/12413 |
| Bug commit | `2d75f539e148` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `markup/goldmark/hugocontext/hugocontext.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00049b180 by goroutine 24:
  runtime.slicecopy()
      /usr/local/go/src/runtime/slice.go:325 +0x0
  bytes.(*Buffer).WriteString()
      /usr/local/go/src/bytes/buffer.go:193 +0x118
  ase/hugo-12413.Wrap()
      /work/hugocontext.go:20 +0xc4
  ase/hugo-12413.TestRace_hugo_12413.func1()
      /work/verified_test.go:20 +0xdc
  ase/hugo-12413.TestRace_hugo_12413.gowrap1()
      /work/verified_test.go:28 +0x41

Previous read at 0x00c00049b184 by goroutine 27:
  ase/hugo-12413.TestRace_hugo_12413.func1()
      /work/verified_test.go:24 +0x10e
  ase/hugo-12413.TestRace_hugo_12413.gowrap1()
      /work/verified_test.go:28 +0x41

Goroutine 24 (running) created at:
  ase/hugo-12413.TestRace_hugo_12413()
      /work/verified_test.go:16 +0x7a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 27 (finished) created at:
  ase/hugo-12413.TestRace_hugo_12413()
      /work/verified_test.go:16 +0x7a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-hugo-12413-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-hugo-12413-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-hugo-12413-fix .
docker run --rm --memory=2g --cpus=1 gonb-hugo-12413-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-hugo-12413-bug .
# (then run as above, no --ssh flag)
```
