# thanos-6575

| Field | Value |
|---|---|
| Project | thanos |
| Reference | https://github.com/thanos-io/thanos/pull/6575 |
| Bug commit | `a35a5b2e519b` |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/store/bucket.go` |


## Race report excerpt

The following stack trace is captured by Go's race detector when running the bug build:

```
WARNING: DATA RACE
Read at 0x00c000048ee8 by goroutine 117:
  github.com/thanos-io/thanos/pkg/store.(*bucketIndexReader).ExpandedPostings.func1()
      /workspace/pkg/store/bucket.go:2208 +0x64
  sort.insertionSort_func()
      /usr/local/go/src/sort/zsortfunc.go:12 +0xc1
  sort.pdqsort_func()
      /usr/local/go/src/sort/zsortfunc.go:73 +0x39e
  sort.Slice()
      /usr/local/go/src/sort/slice.go:29 +0xaa
  github.com/thanos-io/thanos/pkg/store.(*bucketIndexReader).ExpandedPostings()
      /workspace/pkg/store/bucket.go:2207 +0x13c
  github.com/thanos-io/thanos/pkg/store.TestRace_6575.func2()
      /workspace/pkg/store/race_6575_test.go:80 +0x114
  github.com/thanos-io/thanos/pkg/store.TestRace_6575.gowrap1()
      /workspace/pkg/store/race_6575_test.go:82 +0x41

Previous write at 0x00c000048ee8 by goroutine 106:
  internal/reflectlite.Swapper.func3()
      /usr/local/go/src/internal/reflectlite/swapper.go:42 +0xae
  sort.insertionSort_func()
      /usr/local/go/src/sort/zsortfunc.go:13 +0x83
  sort.pdqsort_func()
      /usr/local/go/src/sort/zsortfunc.go:73 +0x39e
  sort.Slice()
      /usr/local/go/src/sort/slice.go:29 +0xaa
  github.com/thanos-io/thanos/pkg/store.(*bucketIndexReader).ExpandedPostings()
      /workspace/pkg/store/bucket.go:2207 +0x13c
  github.com/thanos-io/thanos/pkg/store.TestRace_6575.func2()
      /workspace/pkg/store/race_6575_test.go:80 +0x114
  github.com/thanos-io/thanos/pkg/store.TestRace_6575.gowrap1()
      /workspace/pkg/store/race_6575_test.go:82 +0x41

Goroutine 117 (running) created at:
  github.com/thanos-io/thanos/pkg/store.TestRace_6575()
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
DOCKER_BUILDKIT=1 docker build --ssh default -f bug.Dockerfile -t gonb-thanos-6575-bug .
```

### 3. Trigger race
```bash
docker run --rm --memory=2g --cpus=1 gonb-thanos-6575-bug \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: WARNING: DATA RACE + FAIL
```

### 4. Verify fix
```bash
DOCKER_BUILDKIT=1 docker build --ssh default -f fix.Dockerfile -t gonb-thanos-6575-fix .
docker run --rm --memory=2g --cpus=1 gonb-thanos-6575-fix \
  sh -c "cd /work/pr2t-test && go test -race -vet=off -count=20 -timeout=180s ./..."
# Expected: PASS (race not triggered)
```

## HTTPS fallback (if SSH blocked)

If `git@github.com:` clone fails in your environment:
```bash
sed -i 's|git@github.com:|https://github.com/|g' bug.Dockerfile fix.Dockerfile
# Also remove the --mount=type=ssh hint (HTTPS doesn't need it)
sed -i 's|--mount=type=ssh ||g' bug.Dockerfile fix.Dockerfile
DOCKER_BUILDKIT=1 docker build -f bug.Dockerfile -t gonb-thanos-6575-bug .
# (then run as above, no --ssh flag)
```
