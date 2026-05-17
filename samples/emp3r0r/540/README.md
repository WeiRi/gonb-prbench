# emp3r0r-540

| Field | Value |
|---|---|
| Project | emp3r0r |
| Reference | https://github.com/emp3r0r/emp3r0r/pull/540 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `thudart` |
| Base image | `golang:1.22` (built by gonb-prebuild for this sample) |

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t emp3r0r-540-bug .
docker run --rm --cpus=2 --memory=2g emp3r0r-540-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t emp3r0r-540-fix .
docker run --rm --cpus=2 --memory=2g emp3r0r-540-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
