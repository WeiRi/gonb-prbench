# kubernetes-102090

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/102090 |
| Bug commit | `ec5ec0804d5f` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/client-go/testing/fake.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0002b03a8 by goroutine 15:
  k8s.io/client-go/testing.(*Fake).AddWatchReactor()
      /work/staging/src/k8s.io/client-go/testing/fake.go:109 +0x26b
  k8s.io/client-go/testing.TestRace_102090.func2()
      /work/staging/src/k8s.io/client-go/testing/race_102090_capture_test.go:38 +0xe4
  k8s.io/client-go/testing.TestRace_102090.func4()
      /work/staging/src/k8s.io/client-go/testing/race_102090_capture_test.go:40 +0x41

Previous write at 0x00c0002b03a8 by goroutine 11:
  k8s.io/client-go/testing.(*Fake).AddWatchReactor()
      /work/staging/src/k8s.io/client-go/testing/fake.go:109 +0x32c
  k8s.io/client-go/testing.TestRace_102090.func2()
      /work/staging/src/k8s.io/client-go/testing/race_102090_capture_test.go:38 +0xe4
  k8s.io/client-go/testing.TestRace_102090.func4()
      /work/staging/src/k8s.io/client-go/testing/race_102090_capture_test.go:40 +0x41

Goroutine 15 (running) created at:
  k8s.io/client-go/testing.TestRace_102090()
      /work/staging/src/k8s.io/client-go/testing/race_102090_capture_test.go:35 +0x270
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44

Goroutine 11 (running) created at:
  k8s.io/client-go/testing.TestRace_102090()
      /work/staging/src/k8s.io/client-go/testing/race_102090_capture_test.go:35 +0x270
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1595 +0x261
  testing.(*T).Run.func1()
      /usr/local/go/src/testing/testing.go:1648 +0x44
==================
==================
WARNING: DATA RACE
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-102090-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-102090-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-102090-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-102090-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-102090-bug .
# (then run as above, no --ssh flag)
```
