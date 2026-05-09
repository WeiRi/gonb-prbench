# kubernetes-124344

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/124344 |
| Bug commit | `646fbe6d0a3f` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/client-go/tools/cache/delta_fifo.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000194130 by goroutine 17:
  ase/kubernetes-124344.TestRace_PR124344_TransformerResync.func2()
      /work/verified_test.go:47 +0x4e
  ase/kubernetes-124344.(*DeltaFIFO).Replace()
      /work/delta_fifo.go:45 +0x133
  ase/kubernetes-124344.TestRace_PR124344_TransformerResync.func3()
      /work/verified_test.go:77 +0xf9

Previous write at 0x00c000194130 by goroutine 9:
  ase/kubernetes-124344.TestRace_PR124344_TransformerResync.func2()
      /work/verified_test.go:47 +0x66
  ase/kubernetes-124344.(*DeltaFIFO).Replace()
      /work/delta_fifo.go:45 +0x133
  ase/kubernetes-124344.TestRace_PR124344_TransformerResync.func3()
      /work/verified_test.go:77 +0xf9

Goroutine 17 (running) created at:
  ase/kubernetes-124344.TestRace_PR124344_TransformerResync()
      /work/verified_test.go:70 +0x430
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (running) created at:
  ase/kubernetes-124344.TestRace_PR124344_TransformerResync()
      /work/verified_test.go:70 +0x430
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-124344-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-124344-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-124344-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-124344-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-124344-bug .
# (then run as above, no --ssh flag)
```
