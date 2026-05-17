# kubernetes-134390

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/134390 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/component-base/metrics/metric.go` |
| Base image | `gonb-kubernetes-134390-bug:latest` (built by gonb-prebuild for this sample) |

**Soft issue — source overlay**: `fix.Dockerfile` overlays our hand-written `<src>_fixed.go` instead of applying the original PR `fix.diff`. The PR diff is included for reference but may differ from what `fix.Dockerfile` actually applies. Useful for race reproduction; downstream tools should diff `<src>.go` vs `<src>_fixed.go` for the actual applied fix.

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t kubernetes-134390-bug .
docker run --rm --cpus=2 --memory=2g kubernetes-134390-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t kubernetes-134390-fix .
docker run --rm --cpus=2 --memory=2g kubernetes-134390-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
