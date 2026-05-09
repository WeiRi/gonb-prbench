# kubernetes-133587

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/133587 |
| Bug commit | `17d6c9c551f9` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/dynamic-resource-allocation/structured/internal/experimental/allocator_experimental.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
==================
WARNING: DATA RACE
Read at 0x00c00048c000 by goroutine 62:
  k8s.io/dynamic-resource-allocation/structured/internal/experimental.(*allocator).allocateOne()
      /k8s/staging/src/k8s.io/dynamic-resource-allocation/structured/internal/experimental/allocator_experimental.go:827 +0x246
  k8s.io/dynamic-resource-allocation/structured/internal/experimental.(*Allocator).Allocate()
      /k8s/staging/src/k8s.io/dynamic-resource-allocation/structured/internal/experimental/allocator_experimental.go:342 +0x1a44
  k8s.io/dynamic-resource-allocation/structured/internal/experimental.TestAllocatorClaimsToAllocateRace_133587.func2()
      /k8s/staging/src/k8s.io/dynamic-resource-allocation/structured/internal/experimental/133587_handcrafted_race_test.go:161 +0x108

Previous write at 0x00c00048c000 by goroutine 59:
  k8s.io/dynamic-resource-allocation/structured/internal/experimental.TestAllocatorClaimsToAllocateRace_133587.func2()
      /k8s/staging/src/k8s.io/dynamic-resource-allocation/structured/internal/experimental/133587_handcrafted_race_test.go:161 +0x3c4

Goroutine 62 (running) created at:
  k8s.io/dynamic-resource-allocation/structured/internal/experimental.TestAllocatorClaimsToAllocateRace_133587()
      /k8s/staging/src/k8s.io/dynamic-resource-allocation/structured/internal/experimental/133587_handcrafted_race_test.go:155 +0x7f1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44

Goroutine 59 (running) created at:
  k8s.io/dynamic-resource-allocation/structured/internal/experimental.TestAllocatorClaimsToAllocateRace_133587()
      /k8s/staging/src/k8s.io/dynamic-resource-allocation/structured/internal/experimental/133587_handcrafted_race_test.go:155 +0x7f1
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c000488108 by goroutine 62:
  k8s.io/dynamic-resource-allocation/structured/internal/experimental.(*allocator).allocateOne()
      /k8s/staging/src/k8s.io/dynamic-resource-allocation/structured/internal/experimental/allocator_experimental.go:828 +0x277
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-133587-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-133587-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-133587-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-133587-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-133587-bug .
# (then run as above, no --ssh flag)
```
