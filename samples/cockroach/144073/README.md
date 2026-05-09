# cockroach-144073

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/144073 |
| Bug commit | `6a48f58cf5d7` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/sql/catalog/lease/descriptor_state.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c000012148 by goroutine 11:
  ase/cockroach-144073.(*descriptorVersionState).SetLease()
      /work/descriptor_state.go:49 +0x99
  ase/cockroach-144073.Test144073Race.func1()
      /work/verified_test.go:23 +0xd6
  ase/cockroach-144073.Test144073Race.gowrap1()
      /work/verified_test.go:25 +0x41

Previous read at 0x00c000012148 by goroutine 10:
  ase/cockroach-144073.(*descriptorState).removeInactiveVersions()
      /work/descriptor_state.go:36 +0x1e5
  ase/cockroach-144073.Test144073Race.func1()
      /work/verified_test.go:21 +0xa4
  ase/cockroach-144073.Test144073Race.gowrap1()
      /work/verified_test.go:25 +0x41

Goroutine 11 (running) created at:
  ase/cockroach-144073.Test144073Race()
      /work/verified_test.go:18 +0x20d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 10 (finished) created at:
  ase/cockroach-144073.Test144073Race()
      /work/verified_test.go:18 +0x20d
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-144073-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-144073-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-144073-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-144073-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-144073-bug .
# (then run as above, no --ssh flag)
```
