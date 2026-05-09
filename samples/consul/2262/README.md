# consul-2262

| Field | Value |
|---|---|
| Project | consul |
| Reference | https://github.com/hashicorp/consul/pull/2262 |
| Bug commit | `0047f09be55f` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `command/agent/gated_writer.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000161b8 by goroutine 10:
  ase/consul-2262.(*GatedWriter).Write()
      /work/gated_writer.go:26 +0x155
  ase/consul-2262.TestGatedWriterRace.func1()
      /work/verified_test.go:30 +0xfb

Previous write at 0x00c0000161b8 by goroutine 9:
  ase/consul-2262.(*GatedWriter).Write()
      /work/gated_writer.go:26 +0x215
  ase/consul-2262.TestGatedWriterRace.func1()
      /work/verified_test.go:30 +0xfb

Goroutine 10 (running) created at:
  ase/consul-2262.TestGatedWriterRace()
      /work/verified_test.go:27 +0x139
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44

Goroutine 9 (running) created at:
  ase/consul-2262.TestGatedWriterRace()
      /work/verified_test.go:27 +0x139
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1689 +0x21e
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1742 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c000124ff8 by goroutine 10:
  runtime.growslice()
      /usr/local/go/src/runtime/slice.go:155 +0x0
  ase/consul-2262.(*GatedWriter).Write()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-consul-2262-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-consul-2262-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-consul-2262-fix .
docker run --rm --memory=2g --cpus=1 gonb-consul-2262-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-consul-2262-bug .
# (then run as above, no --ssh flag)
```
