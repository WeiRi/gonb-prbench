# nomad-14188

| Field | Value |
|---|---|
| Project | nomad |
| Reference | https://github.com/nomad/nomad/pull/14188 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `nomad/stream/event_broker.go` |
| Base image | `gonb-nomad-14188-base:latest` (built by gonb-prebuild for this sample) |

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t nomad-14188-bug .
docker run --rm --cpus=2 --memory=2g nomad-14188-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t nomad-14188-fix .
docker run --rm --cpus=2 --memory=2g nomad-14188-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
