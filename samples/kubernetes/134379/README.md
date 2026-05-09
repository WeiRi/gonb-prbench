# kubernetes-134379

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/134379 |
| Bug commit | `8ac5701d3a14` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/controller/garbagecollector/graph.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c0000f60d0 by goroutine 14:
  reflect.typedmemmove()
      /usr/local/go/src/runtime/mbarrier.go:203 +0x0
  reflect.packEface()
      /usr/local/go/src/reflect/value.go:135 +0xc5
  reflect.valueInterface()
      /usr/local/go/src/reflect/value.go:1526 +0x179
  reflect.Value.Interface()
      /usr/local/go/src/reflect/value.go:1496 +0xb4
  fmt.(*pp).printValue()
      /usr/local/go/src/fmt/print.go:769 +0xc5
  fmt.(*pp).printValue()
      /usr/local/go/src/fmt/print.go:921 +0x132a
  fmt.(*pp).printArg()
      /usr/local/go/src/fmt/print.go:759 +0xb84
  fmt.(*pp).doPrintf()
      /usr/local/go/src/fmt/print.go:1075 +0x592
  fmt.Sprintf()
      /usr/local/go/src/fmt/print.go:239 +0x5c
  ase/kubernetes-134379.(*node).String()
      /work/graph.go:33 +0xfd
  ase/kubernetes-134379.TestRace_134379.func1()
      /work/verified_test.go:27 +0xa4

Previous write at 0x00c0000f60d0 by goroutine 8:
  ase/kubernetes-134379.(*node).markBeingDeleted()
      /work/graph.go:40 +0xc5
  ase/kubernetes-134379.TestRace_134379.func1()
      /work/verified_test.go:28 +0xbd

Goroutine 14 (running) created at:
  ase/kubernetes-134379.TestRace_134379()
      /work/verified_test.go:24 +0x1ec
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-134379-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-134379-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-134379-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-134379-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-134379-bug .
# (then run as above, no --ssh flag)
```
