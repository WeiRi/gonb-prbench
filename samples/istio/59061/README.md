# istio-59061

| Field | Value |
|---|---|
| Project | istio |
| Reference | https://github.com/istio/istio/pull/59061 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/kube/multicluster/cluster.go` |
| Base image | `inp-istio-59061` (built by gonb-prebuild for this sample) |

**Soft issue — dual-test**: `fix.Dockerfile` swaps in a different test file (`verified_test_fixed.go` or `verified_test_inplace.go`) because PR's fix changed an API that the original test referenced. Both bug & fix test the same race, but via different test code.

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t istio-59061-bug .
docker run --rm --cpus=2 --memory=2g istio-59061-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t istio-59061-fix .
docker run --rm --cpus=2 --memory=2g istio-59061-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
