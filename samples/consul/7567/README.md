# consul-7567

| Field | Value |
|---|---|
| Project | consul |
| Reference | https://github.com/hashicorp/consul/pull/7567 |
| Bug commit | `fce56e4cb68c` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `sdk/freeport/freeport.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0000d4598 by goroutine 14:
  container/list.(*List).insert()
      /usr/local/go/src/container/list/list.go:98 +0x7c
  container/list.(*List).insertValue()
      /usr/local/go/src/container/list/list.go:104 +0x1ac
  container/list.(*List).PushBack()
      /usr/local/go/src/container/list/list.go:152 +0x192
  ase/consul-7567.Return()
      /work/freeport.go:60 +0xd2
  ase/consul-7567.TestFreeportResetRace.func1()
      /work/verified_test.go:32 +0x89

Previous read at 0x00c0000d4598 by goroutine 8:
  container/list.(*List).Front()
      /usr/local/go/src/container/list/list.go:70 +0x67
  ase/consul-7567.checkFreedPorts()
      /work/freeport.go:70 +0x42

Goroutine 14 (running) created at:
  ase/consul-7567.TestFreeportResetRace()
      /work/verified_test.go:26 +0x249
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 8 (running) created at:
  ase/consul-7567.initialize()
      /work/freeport.go:28 +0x558
  ase/consul-7567.TestFreeportResetRace()
      /work/verified_test.go:19 +0x67
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-consul-7567-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-consul-7567-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-consul-7567-fix .
docker run --rm --memory=2g --cpus=1 gonb-consul-7567-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-consul-7567-bug .
# (then run as above, no --ssh flag)
```
