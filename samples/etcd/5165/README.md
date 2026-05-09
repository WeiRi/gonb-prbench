# etcd-5165

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd-io/etcd/pull/5165 |
| Bug commit | `0dd9c2520b69` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `raft/node.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000192040 by goroutine 58:
  ase/etcd-5165.(*node).ReadPropc()
      /work/node.go:26 +0x84
  ase/etcd-5165.TestRace_5165.func1.2()
      /work/verified_test.go:37 +0x7d

Previous write at 0x00c000192040 by goroutine 57:
  ase/etcd-5165.(*node).ClearPropc()
      /work/node.go:21 +0x84
  ase/etcd-5165.TestRace_5165.func1.1()
      /work/verified_test.go:33 +0x7d

Goroutine 58 (running) created at:
  ase/etcd-5165.TestRace_5165.func1()
      /work/verified_test.go:35 +0xaf

Goroutine 57 (finished) created at:
  ase/etcd-5165.TestRace_5165.func1()
      /work/verified_test.go:31 +0x4a4
==================
    testing.go:1398: race detected during execution of test
--- FAIL: TestRace_5165 (0.70s)
FAIL
FAIL	ase/etcd-5165	0.731s
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-etcd-5165-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-etcd-5165-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-etcd-5165-fix .
docker run --rm --memory=2g --cpus=1 gonb-etcd-5165-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-etcd-5165-bug .
# (then run as above, no --ssh flag)
```
