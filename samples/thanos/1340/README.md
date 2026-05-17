# thanos-1340

| Field | Value |
|---|---|
| Project | thanos |
| Reference | https://github.com/thanos/thanos/pull/1340 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/pool/pool.go` |
| Base image | `gonb-thanos-1340-bug:latest` (built by gonb-prebuild for this sample) |

**Soft issue — source overlay**: `fix.Dockerfile` overlays our hand-written `<src>_fixed.go` instead of applying the original PR `fix.diff`. The PR diff is included for reference but may differ from what `fix.Dockerfile` actually applies. Useful for race reproduction; downstream tools should diff `<src>.go` vs `<src>_fixed.go` for the actual applied fix.

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t thanos-1340-bug .
docker run --rm --cpus=2 --memory=2g thanos-1340-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t thanos-1340-fix .
docker run --rm --cpus=2 --memory=2g thanos-1340-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
