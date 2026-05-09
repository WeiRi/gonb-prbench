# kubernetes-133781

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/133781 |
| Bug commit | `947a8ebfd14f` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/apis/scheduling/v1/helpers.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x000000716c80 by goroutine 13:
  ase/kubernetes-133781.IsKnownSystemPriorityClass()
      /work/helpers.go:28 +0x136
  ase/kubernetes-133781.TestRace_133781.func2()
      /work/verified_test.go:29 +0xb9

Previous write at 0x000000716c80 by goroutine 8:
  ase/kubernetes-133781.TestRace_133781.func1()
      /work/verified_test.go:21 +0x12f

Goroutine 13 (running) created at:
  ase/kubernetes-133781.TestRace_133781()
      /work/verified_test.go:26 +0x6f
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (finished) created at:
  ase/kubernetes-133781.TestRace_133781()
      /work/verified_test.go:16 +0x11b
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Write at 0x000000716cc0 by goroutine 12:
  ase/kubernetes-133781.TestRace_133781.func1()
      /work/verified_test.go:21 +0x12f

Previous write at 0x000000716cc0 by goroutine 8:
  ase/kubernetes-133781.TestRace_133781.func1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-133781-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-133781-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-133781-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-133781-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-133781-bug .
# (then run as above, no --ssh flag)
```
