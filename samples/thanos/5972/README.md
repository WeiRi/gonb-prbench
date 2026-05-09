# thanos-5972

| Field | Value |
|---|---|
| Project | thanos |
| Reference | https://github.com/thanos-io/thanos/pull/5972 |
| Bug commit | `d76c723161df` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/query/endpointset.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Write at 0x00c0005188e8 by goroutine 139:
  github.com/thanos-io/thanos/pkg/query.TestRace_5972.func4()
      /workspace/pkg/query/race_5972_test.go:75 +0x16e

Previous read at 0x00c0005188e8 by goroutine 115:
  github.com/thanos-io/thanos/pkg/query.(*EndpointSet).getTimedOutRefs()
      /workspace/pkg/query/endpointset.go:455 +0x4b
  github.com/thanos-io/thanos/pkg/query.TestRace_5972.func3()
      /workspace/pkg/query/race_5972_test.go:63 +0xa4

Goroutine 139 (running) created at:
  github.com/thanos-io/thanos/pkg/query.TestRace_5972()
      /workspace/pkg/query/race_5972_test.go:71 +0x87c
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44

Goroutine 115 (running) created at:
  github.com/thanos-io/thanos/pkg/query.TestRace_5972()
      /workspace/pkg/query/race_5972_test.go:60 +0x647
  testing.tRunner()
      /usr/local/go/src/testing/testing.go:1934 +0x21c
  testing.(*T).Run.gowrap1()
      /usr/local/go/src/testing/testing.go:1997 +0x44
==================
==================
WARNING: DATA RACE
Read at 0x00c000195e00 by goroutine 124:
  runtime.mapIterStart()
      /usr/local/go/src/runtime/map_swiss.go:160 +0x0
  github.com/thanos-io/thanos/pkg/query.(*EndpointSet).getTimedOutRefs()
      /workspace/pkg/query/endpointset.go:457 +0xc8
  github.com/thanos-io/thanos/pkg/query.TestRace_5972.func3()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-thanos-5972-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-thanos-5972-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-thanos-5972-fix .
docker run --rm --memory=2g --cpus=1 gonb-thanos-5972-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-thanos-5972-bug .
# (then run as above, no --ssh flag)
```
