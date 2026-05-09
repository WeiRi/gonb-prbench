# thanos-5503

| Field | Value |
|---|---|
| Project | thanos |
| Reference | https://github.com/thanos-io/thanos/pull/5503 |
| Bug commit | `9b6903b58c23` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/compact/compact.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0005381f0 by goroutine 485:
  github.com/thanos-io/thanos/pkg/compact.(*Group).compact.(*Group).compact.func2.func7.1()
      /workspace/pkg/compact/compact.go:1033 +0x1fb
  github.com/thanos-io/thanos/pkg/compact.(*Group).compact.(*Group).compact.func2.func7.1()
      /workspace/pkg/compact/compact.go:1033 +0x1de
  github.com/thanos-io/thanos/pkg/tracing.DoInSpanWithErr()
      /workspace/pkg/tracing/tracing.go:82 +0x14f
  github.com/thanos-io/thanos/pkg/compact.(*Group).compact.(*Group).compact.func2.func7()
      /workspace/pkg/compact/compact.go:1032 +0x270
  golang.org/x/sync/errgroup.(*Group).Go.func1()
      /go/pkg/mod/golang.org/x/sync@v0.0.0-20220601150217-0de741cfad7f/errgroup/errgroup.go:75 +0x86

Previous write at 0x00c0005381f0 by goroutine 486:
  github.com/thanos-io/thanos/pkg/compact.(*Group).compact.(*Group).compact.func2.func7.1()
      /workspace/pkg/compact/compact.go:1033 +0x1fb
  github.com/thanos-io/thanos/pkg/compact.(*Group).compact.(*Group).compact.func2.func7.1()
      /workspace/pkg/compact/compact.go:1033 +0x1de
  github.com/thanos-io/thanos/pkg/tracing.DoInSpanWithErr()
      /workspace/pkg/tracing/tracing.go:82 +0x14f
  github.com/thanos-io/thanos/pkg/compact.(*Group).compact.(*Group).compact.func2.func7()
      /workspace/pkg/compact/compact.go:1032 +0x270
  golang.org/x/sync/errgroup.(*Group).Go.func1()
      /go/pkg/mod/golang.org/x/sync@v0.0.0-20220601150217-0de741cfad7f/errgroup/errgroup.go:75 +0x86

Goroutine 485 (running) created at:
  golang.org/x/sync/errgroup.(*Group).Go()
      /go/pkg/mod/golang.org/x/sync@v0.0.0-20220601150217-0de741cfad7f/errgroup/errgroup.go:72 +0x11c
  github.com/thanos-io/thanos/pkg/compact.(*Group).compact.func2()
      /workspace/pkg/compact/compact.go:1031 +0xe86
  github.com/thanos-io/thanos/pkg/compact.(*Group).compact()
      /workspace/pkg/compact/compact.go:1068 +0xcda
  github.com/thanos-io/thanos/pkg/compact.(*Group).Compact.func2()
      /workspace/pkg/compact/compact.go:778 +0xfb
  github.com/thanos-io/thanos/pkg/tracing.DoInSpanWithErr()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-thanos-5503-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-thanos-5503-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-thanos-5503-fix .
docker run --rm --memory=2g --cpus=1 gonb-thanos-5503-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-thanos-5503-bug .
# (then run as above, no --ssh flag)
```
