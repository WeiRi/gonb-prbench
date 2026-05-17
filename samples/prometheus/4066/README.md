# prometheus-4066

| Field | Value |
|---|---|
| Project | prometheus |
| Reference | https://github.com/prometheus/prometheus/pull/4066 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `scrape/target.go` |
| Base image | `d3-prom-4066-base:latest` (built by gonb-prebuild for this sample) |

**Soft issue — dual-test**: `fix.Dockerfile` swaps in a different test file (`verified_test_fixed.go` or `verified_test_inplace.go`) because PR's fix changed an API that the original test referenced. Both bug & fix test the same race, but via different test code.

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t prometheus-4066-bug .
docker run --rm --cpus=2 --memory=2g prometheus-4066-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t prometheus-4066-fix .
docker run --rm --cpus=2 --memory=2g prometheus-4066-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
