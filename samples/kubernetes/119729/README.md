# kubernetes-119729

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/119729 |
| Bug commit | `99190634ab25` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/scheduler/schedule_one.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000188000 by goroutine 20:
  ase/kubernetes-119729.(*scheduler).writePod()
      /work/schedule_one.go:25 +0x10d
  ase/kubernetes-119729.(*scheduler).scheduleOne.func1()
      /work/schedule_one.go:40 +0x12

Previous read at 0x00c000188000 by goroutine 8:
  ase/kubernetes-119729.(*scheduler).done()
      /work/schedule_one.go:30 +0xf1
  ase/kubernetes-119729.(*scheduler).scheduleOne()
      /work/schedule_one.go:42 +0xec
  ase/kubernetes-119729.TestScheduleOneBindingFailureRace.func1()
      /work/verified_test.go:31 +0x14f

Goroutine 20 (running) created at:
  ase/kubernetes-119729.(*scheduler).scheduleOne()
      /work/schedule_one.go:40 +0xe4
  ase/kubernetes-119729.TestScheduleOneBindingFailureRace.func1()
      /work/verified_test.go:31 +0x14f

Goroutine 8 (finished) created at:
  ase/kubernetes-119729.TestScheduleOneBindingFailureRace()
      /work/verified_test.go:28 +0x84
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestScheduleOneBindingFailureRace (0.00s)
FAIL
FAIL	ase/kubernetes-119729	0.021s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-119729-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-119729-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-119729-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-119729-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-119729-bug .
# (then run as above, no --ssh flag)
```
