# moby-26695

| Field | Value |
|---|---|
| Project | moby |
| Reference | https://github.com/moby/moby/pull/26695 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `libcontainerd/pausemonitor_linux.go` |
| Base image | `gonb-moby-26695-base:latest` (built by gonb-prebuild for this sample) |

**Soft issue — source overlay**: `fix.Dockerfile` overlays our hand-written `<src>_fixed.go` instead of applying the original PR `fix.diff`. The PR diff is included for reference but may differ from what `fix.Dockerfile` actually applies. Useful for race reproduction; downstream tools should diff `<src>.go` vs `<src>_fixed.go` for the actual applied fix.

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t moby-26695-bug .
docker run --rm --cpus=2 --memory=2g moby-26695-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t moby-26695-fix .
docker run --rm --cpus=2 --memory=2g moby-26695-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
