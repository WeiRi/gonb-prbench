# swarmkit-1046

| Field | Value |
|---|---|
| Project | swarmkit |
| Reference | https://github.com/moby/swarmkit/pull/1046 |
| Bug commit | `5f1106a1fd8b` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `agent/node.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000114008 by goroutine 33:
  ase/swarmkit-1046.RunManagerLoopBUG.func1()
      /work/node.go:17 +0x5b

Previous write at 0x00c000114008 by goroutine 8:
  ase/swarmkit-1046.RunManagerLoopBUG()
      /work/node.go:25 +0x1c6
  ase/swarmkit-1046.TestRace_swarmkit_1046.func1()
      /work/verified_test.go:21 +0x104

Goroutine 33 (running) created at:
  ase/swarmkit-1046.RunManagerLoopBUG()
      /work/node.go:14 +0x150
  ase/swarmkit-1046.TestRace_swarmkit_1046.func1()
      /work/verified_test.go:21 +0x104

Goroutine 8 (finished) created at:
  ase/swarmkit-1046.TestRace_swarmkit_1046()
      /work/verified_test.go:15 +0x64
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_swarmkit_1046 (0.00s)
FAIL
FAIL	ase/swarmkit-1046	0.022s
FAIL
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-swarmkit-1046-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-swarmkit-1046-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-swarmkit-1046-fix .
docker run --rm --memory=2g --cpus=1 gonb-swarmkit-1046-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-swarmkit-1046-bug .
# (then run as above, no --ssh flag)
```
