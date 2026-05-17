# swarmkit-1046

| Field | Value |
|---|---|
| Project | swarmkit |
| Reference | https://github.com/swarmkit/swarmkit/pull/1046 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `agent/node.go` |
| Base image | `gonb-swarmkit-1046-base-v6:latest` (built by gonb-prebuild for this sample) |

**Soft issue — H9 drift**: race detector frame is NOT in any file modified by `fix.diff`. The fix likely suppresses race via a side-effect path (e.g. removes a goroutine that walked into the racy code) rather than directly addressing the racing line. Effective for bug reproduction; downstream fix-experiment tools should verify root-cause is addressed.

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t swarmkit-1046-bug .
docker run --rm --cpus=2 --memory=2g swarmkit-1046-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t swarmkit-1046-fix .
docker run --rm --cpus=2 --memory=2g swarmkit-1046-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
