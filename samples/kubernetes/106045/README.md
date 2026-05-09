# kubernetes-106045

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/106045 |
| Bug commit | `c04157895ca7` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `audit.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00011e3f0 by goroutine 9:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/kubernetes-106045.(*auditHandler).logAnnotations()
      /work/audit.go:24 +0x1a8
  ase/kubernetes-106045.(*auditHandler).Admit()
      /work/audit.go:29 +0x1ee
  ase/kubernetes-106045.TestRace_106045_AuditAnnotations.func2()
      /work/verified_test.go:24 +0xcf

Previous write at 0x00c00011e3f0 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/go/src/runtime/map_faststr.go:203 +0x0
  ase/kubernetes-106045.(*auditHandler).logAnnotations()
      /work/audit.go:24 +0x1a8
  ase/kubernetes-106045.(*auditHandler).Admit()
      /work/audit.go:29 +0x1ee
  ase/kubernetes-106045.TestRace_106045_AuditAnnotations.func1()
      /work/verified_test.go:18 +0xcf

Goroutine 9 (running) created at:
  ase/kubernetes-106045.TestRace_106045_AuditAnnotations()
      /work/verified_test.go:21 +0x56
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 8 (finished) created at:
  ase/kubernetes-106045.TestRace_106045_AuditAnnotations()
      /work/verified_test.go:15 +0x1dc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-106045-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-106045-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-106045-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-106045-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-106045-bug .
# (then run as above, no --ssh flag)
```
