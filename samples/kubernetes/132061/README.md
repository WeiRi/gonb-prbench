# kubernetes-132061

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/132061 |
| Bug commit | `fd53f7292c7d` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/apiserver/pkg/cel/common/typeprovider.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
==================
WARNING: DATA RACE
Write at 0x00c00030aa50 by goroutine 61:
  k8s.io/apiserver/pkg/cel/common.TestResolverEnvOptionRace_132061.ResolverEnvOption.NewResolverTypeProviderAndEnvOption.func2()
      /k8s/staging/src/k8s.io/apiserver/pkg/cel/common/typeprovider.go:122 +0x4d
  github.com/google/cel-go/cel.(*Env).configure()
      /k8s/vendor/github.com/google/cel-go/cel/env.go:752 +0xa7
  github.com/google/cel-go/cel.(*Env).Extend()
      /k8s/vendor/github.com/google/cel-go/cel/env.go:525 +0x13d4
  github.com/google/cel-go/cel.NewEnv()
      /k8s/vendor/github.com/google/cel-go/cel/env.go:304 +0x1cc
  k8s.io/apiserver/pkg/cel/common.TestResolverEnvOptionRace_132061.func1()
      /k8s/staging/src/k8s.io/apiserver/pkg/cel/common/132061_handcrafted_race_test.go:50 +0xc5

Previous write at 0x00c00030aa50 by goroutine 59:
  k8s.io/apiserver/pkg/cel/common.TestResolverEnvOptionRace_132061.ResolverEnvOption.NewResolverTypeProviderAndEnvOption.func2()
      /k8s/staging/src/k8s.io/apiserver/pkg/cel/common/typeprovider.go:122 +0x4d
  github.com/google/cel-go/cel.(*Env).configure()
      /k8s/vendor/github.com/google/cel-go/cel/env.go:752 +0xa7
  github.com/google/cel-go/cel.(*Env).Extend()
      /k8s/vendor/github.com/google/cel-go/cel/env.go:525 +0x13d4
  github.com/google/cel-go/cel.NewEnv()
      /k8s/vendor/github.com/google/cel-go/cel/env.go:304 +0x1cc
  k8s.io/apiserver/pkg/cel/common.TestResolverEnvOptionRace_132061.func1()
      /k8s/staging/src/k8s.io/apiserver/pkg/cel/common/132061_handcrafted_race_test.go:50 +0xc5

Goroutine 61 (running) created at:
  k8s.io/apiserver/pkg/cel/common.TestResolverEnvOptionRace_132061()
      /k8s/staging/src/k8s.io/apiserver/pkg/cel/common/132061_handcrafted_race_test.go:47 +0x15d
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1792 +0x225
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1851 +0x44

Goroutine 59 (running) created at:
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-kubernetes-132061-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-132061-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-kubernetes-132061-fix .
docker run --rm --memory=2g --cpus=1 gonb-kubernetes-132061-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-kubernetes-132061-bug .
# (then run as above, no --ssh flag)
```
