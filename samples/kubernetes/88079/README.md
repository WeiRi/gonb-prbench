# kubernetes-88079

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/88079 |
| Bug commit | `f7eafa1a838c` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/client-go/util/connrotation/connrotation.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000218370 by goroutine 35:
  ase/kubernetes-88079.(*Dialer).CloseAll()
      /work/connrotation.go:64 +0x16f
  ase/kubernetes-88079.TestRace_88079.func3()
      /work/verified_test.go:41 +0xa4

Previous write at 0x00c000218370 by goroutine 16:
  ase/kubernetes-88079.(*Dialer).DialContext()
      /work/connrotation.go:49 +0x244
  ase/kubernetes-88079.(*Dialer).Dial()
      /work/connrotation.go:36 +0xf8
  ase/kubernetes-88079.TestRace_88079.func2()
      /work/verified_test.go:33 +0xc7

Goroutine 35 (running) created at:
  ase/kubernetes-88079.TestRace_88079()
      /work/verified_test.go:38 +0x12a
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 16 (running) created at:
  ase/kubernetes-88079.TestRace_88079()
      /work/verified_test.go:30 +0x21c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_88079 (0.16s)
FAIL
FAIL	ase/kubernetes-88079	0.181s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-88079-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-88079-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-88079-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-88079-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-88079-bug .
# (then run as above, no --ssh flag)
```
