# hugo-7393

| Field | Value |
|---|---|
| Project | hugo |
| Reference | https://github.com/gohugoio/hugo/pull/7393 |
| Bug commit | `522ba1cd98ac` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `hugolib/content_map_page.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000098078 by goroutine 13:
  ase/hugo-7393.(*HugoSites).GetContentPage()
      /work/hugolib_site.go:21 +0x84
  ase/hugo-7393.TestRaceContentField.func2()
      /work/verified_test.go:14 +0x12

Previous write at 0x00c000098078 by goroutine 8:
  ase/hugo-7393.(*HugoSites).readAndProcessContent()
      /work/hugolib_site.go:16 +0xb0
  ase/hugo-7393.TestRaceContentField.func1()
      /work/verified_test.go:13 +0x12

Goroutine 13 (running) created at:
  ase/hugo-7393.TestRaceContentField()
      /work/verified_test.go:14 +0xc4
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/hugo-7393.TestRaceContentField()
      /work/verified_test.go:13 +0x19c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c000098078 by goroutine 12:
  ase/hugo-7393.(*HugoSites).readAndProcessContent()
      /work/hugolib_site.go:15 +0x84
  ase/hugo-7393.TestRaceContentField.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-hugo-7393-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-hugo-7393-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-hugo-7393-fix .
docker run --rm --memory=2g --cpus=1 gonb-hugo-7393-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-hugo-7393-bug .
# (then run as above, no --ssh flag)
```
