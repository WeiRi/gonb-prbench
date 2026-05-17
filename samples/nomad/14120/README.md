# nomad-14120

| Field | Value |
|---|---|
| Project | nomad |
| Reference | https://github.com/nomad/nomad/pull/14120 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `nomad/serf.go` |
| Base image | `gonb-nomad-14120-base:latest` (built by gonb-prebuild for this sample) |

**Soft issue — H9 drift**: race detector frame is NOT in any file modified by `fix.diff`. The fix likely suppresses race via a side-effect path (e.g. removes a goroutine that walked into the racy code) rather than directly addressing the racing line. Effective for bug reproduction; downstream fix-experiment tools should verify root-cause is addressed.

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t nomad-14120-bug .
docker run --rm --cpus=2 --memory=2g nomad-14120-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t nomad-14120-fix .
docker run --rm --cpus=2 --memory=2g nomad-14120-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
