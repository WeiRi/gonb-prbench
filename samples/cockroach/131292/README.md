# cockroach-131292

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/131292 |
| Bug commit | `4882ef1c5dba` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/util/admission/work_queue.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000194120 by goroutine 8:
  ase/cockroach-131292.(*heap).Push()
      /work/work_queue.go:12 +0x167
  ase/cockroach-131292.(*WorkQueue).Push()
      /work/work_queue.go:37 +0x9b
  ase/cockroach-131292.Test131292Race.func1()
      /work/verified_test.go:16 +0xa4
  ase/cockroach-131292.Test131292Race.gowrap1()
      /work/verified_test.go:21 +0x41

Previous read at 0x00c000194120 by goroutine 39:
  ase/cockroach-131292.(*heap).Len()
      /work/work_queue.go:11 +0xe8
  ase/cockroach-131292.(*WorkQueue).Admit()
      /work/work_queue.go:25 +0xa7
  ase/cockroach-131292.Test131292Race.func1()
      /work/verified_test.go:18 +0xb7
  ase/cockroach-131292.Test131292Race.gowrap1()
      /work/verified_test.go:21 +0x41

Goroutine 8 (running) created at:
  ase/cockroach-131292.Test131292Race()
      /work/verified_test.go:13 +0x184
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 39 (finished) created at:
  ase/cockroach-131292.Test131292Race()
      /work/verified_test.go:13 +0x184
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-131292-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-131292-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-131292-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-131292-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-131292-bug .
# (then run as above, no --ssh flag)
```
