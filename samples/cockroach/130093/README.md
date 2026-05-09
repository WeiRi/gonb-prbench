# cockroach-130093

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/130093 |
| Bug commit | `8d6b63f0645d` |
| Category | special_library |
| Oracle | RACE |
| Primary diff file | `pkg/server/license/enforcer.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00000071ffc0 by goroutine 8:
  ase/cockroach-130093.GetEnforcerInstance.func1()
      /work/enforcer.go:25 +0xb1
  sync.(*Once).doSlow()
      /usr/local/go/src/sync/once.go:74 +0xf0
  sync.(*Once).Do()
      /usr/local/go/src/sync/once.go:65 +0x44
  ase/cockroach-130093.GetEnforcerInstance()
      /work/enforcer.go:24 +0x44
  ase/cockroach-130093.Test130092Race.func1()
      /work/verified_test.go:18 +0x76
  ase/cockroach-130093.Test130092Race.gowrap1()
      /work/verified_test.go:24 +0x41

Previous read at 0x00000071ffc0 by goroutine 9:
  ase/cockroach-130093.GetEnforcerInstance()
      /work/enforcer.go:23 +0x24
  ase/cockroach-130093.Test130092Race.func1()
      /work/verified_test.go:18 +0x76
  ase/cockroach-130093.Test130092Race.gowrap1()
      /work/verified_test.go:24 +0x41

Goroutine 8 (running) created at:
  ase/cockroach-130093.Test130092Race()
      /work/verified_test.go:16 +0xc5
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (running) created at:
  ase/cockroach-130093.Test130092Race()
      /work/verified_test.go:16 +0xc5
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-130093-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-130093-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-130093-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-130093-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-130093-bug .
# (then run as above, no --ssh flag)
```
