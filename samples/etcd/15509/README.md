# etcd-15509

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/15509 |
| Bug commit | `736c89398bcf` |
| Category | anonymous_function |
| Oracle | RACE |
| Primary diff file | `server/embed/serve.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000188008 by goroutine 39:
  ase/etcd-15509.runServe.func1()
      /work/serve.go:20 +0x2e

Previous write at 0x00c000188008 by goroutine 8:
  ase/etcd-15509.runServe()
      /work/serve.go:23 +0x1b0
  ase/etcd-15509.TestRace_PR15509_GsClosureCapture.gowrap1()
      /work/verified_test.go:16 +0x33

Goroutine 39 (running) created at:
  ase/etcd-15509.runServe()
      /work/serve.go:19 +0x184
  ase/etcd-15509.TestRace_PR15509_GsClosureCapture.gowrap1()
      /work/verified_test.go:16 +0x33

Goroutine 8 (finished) created at:
  ase/etcd-15509.TestRace_PR15509_GsClosureCapture()
      /work/verified_test.go:16 +0x64
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c00018a048 by goroutine 39:
  ase/etcd-15509.runServe.func1()
      /work/serve.go:20 +0x3d

Previous write at 0x00c00018a048 by goroutine 8:
  ase/etcd-15509.runServe()
      /work/serve.go:23 +0x19a
  ase/etcd-15509.TestRace_PR15509_GsClosureCapture.gowrap1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-15509-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-15509-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-15509-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-15509-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-15509-bug .
# (then run as above, no --ssh flag)
```
