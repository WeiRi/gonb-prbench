# cockroach-67006

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/67006 |
| Bug commit | `bb89f0136239` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kv/kvserver/scanner.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c00009a078 by goroutine 8:
  ase/cockroach-67006.(*replicaScanner).Start()
      /work/scanner.go:17 +0xa7
  ase/cockroach-67006.Test67006Race.func1()
      /work/verified_test.go:14 +0x7a

Previous read at 0x00c00009a078 by goroutine 9:
  ase/cockroach-67006.(*replicaScanner).Monitor()
      /work/scanner.go:30 +0x84
  ase/cockroach-67006.Test67006Race.func2()
      /work/verified_test.go:19 +0x7d

Goroutine 8 (running) created at:
  ase/cockroach-67006.Test67006Race()
      /work/verified_test.go:12 +0x13c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (finished) created at:
  ase/cockroach-67006.Test67006Race()
      /work/verified_test.go:17 +0x1fc
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
--- FAIL: Test67006Race (0.00s)
    testing.go:1398: race detected during execution of test
FAIL
FAIL	ase/cockroach-67006	0.019s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-cockroach-67006-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-cockroach-67006-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-cockroach-67006-fix .
docker run --rm --memory=2g --cpus=1 gonb-cockroach-67006-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-cockroach-67006-bug .
# (then run as above, no --ssh flag)
```
