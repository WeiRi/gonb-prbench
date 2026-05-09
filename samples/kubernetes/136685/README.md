# kubernetes-136685

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/136685 |
| Bug commit | `b23802b6092b` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/apis/rbac/v1/helpers.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c00018c010 by goroutine 10:
  slices.insertionSortOrdered[go.shape.string]()
      /usr/local/go/src/slices/zsortordered.go:14 +0x12a
  slices.pdqsortOrdered[go.shape.string]()
      /usr/local/go/src/slices/zsortordered.go:75 +0x6ef
  slices.Sort[go.shape.[]string,go.shape.string]()
      /usr/local/go/src/slices/sort.go:18 +0x64
  sort.stringsImpl()
      /usr/local/go/src/sort/sort_impl_go121.go:18 +0xe
  sort.Strings()
      /usr/local/go/src/sort/sort.go:176 +0x18b
  ase/kubernetes-136685.(*PolicyRuleBuilder).Rule()
      /work/helpers.go:38 +0x162
  ase/kubernetes-136685.TestRace_136685_SharedVerbsSlice.func1()
      /work/verified_test.go:21 +0x161

Previous write at 0x00c00018c010 by goroutine 9:
  slices.insertionSortOrdered[go.shape.string]()
      /usr/local/go/src/slices/zsortordered.go:15 +0x308
  slices.pdqsortOrdered[go.shape.string]()
      /usr/local/go/src/slices/zsortordered.go:75 +0x6ef
  slices.Sort[go.shape.[]string,go.shape.string]()
      /usr/local/go/src/slices/sort.go:18 +0x64
  sort.stringsImpl()
      /usr/local/go/src/sort/sort_impl_go121.go:18 +0xe
  sort.Strings()
      /usr/local/go/src/sort/sort.go:176 +0x18b
  ase/kubernetes-136685.(*PolicyRuleBuilder).Rule()
      /work/helpers.go:38 +0x162
  ase/kubernetes-136685.TestRace_136685_SharedVerbsSlice.func1()
      /work/verified_test.go:21 +0x161

Goroutine 10 (running) created at:
  ase/kubernetes-136685.TestRace_136685_SharedVerbsSlice()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-136685-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-136685-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-136685-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-136685-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-136685-bug .
# (then run as above, no --ssh flag)
```
