# serving-3068

| Field | Value |
|---|---|
| Project | serving |
| Reference | https://github.com/knative/serving/pull/3068 |
| Bug commit | `899424ba28c7` |
| Category | channel_misuse |
| Oracle | RACE|PANIC |
| Primary diff file | `pkg/pool/pool.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0001b2130 by goroutine 9:
  runtime.chansend()
      /usr/local/go/src/runtime/chan.go:160 +0x0
  ase/serving-3068.(*impl).Go()
      /work/pool.go:50 +0x64
  ase/serving-3068.TestRace_serving3068.func1()
      /work/verified_test.go:38 +0x1c1

Previous write at 0x00c0001b2130 by goroutine 10:
  runtime.closechan()
      /usr/local/go/src/runtime/chan.go:357 +0x0
  ase/serving-3068.(*impl).Wait.func1()
      /work/pool.go:55 +0x45
  sync.(*Once).doSlow()
      /usr/local/go/src/sync/once.go:74 +0xf0
  sync.(*Once).Do()
      /usr/local/go/src/sync/once.go:65 +0x44
  ase/serving-3068.(*impl).Wait()
      /work/pool.go:54 +0x67
  ase/serving-3068.TestRace_serving3068.func2()
      /work/verified_test.go:43 +0x39

Goroutine 9 (running) created at:
  ase/serving-3068.TestRace_serving3068()
      /work/verified_test.go:23 +0x2dd
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 10 (finished) created at:
  ase/serving-3068.TestRace_serving3068()
      /work/verified_test.go:42 +0x34a
  testing.tRunner()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-serving-3068-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-serving-3068-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-serving-3068-fix .
docker run --rm --memory=2g --cpus=1 gonb-serving-3068-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-serving-3068-bug .
# (then run as above, no --ssh flag)
```
