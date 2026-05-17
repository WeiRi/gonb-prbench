# cockroach-130094

| Field | Value |
|---|---|
| Project | cockroach |
| Reference | https://github.com/cockroachdb/cockroach/pull/130094 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `pkg/server/license/enforcer.go` |
| Base image | `gonb-cockroach-130094-base:latest` (built by gonb-prebuild for this sample) |

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t cockroach-130094-bug .
docker run --rm --cpus=2 --memory=2g cockroach-130094-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t cockroach-130094-fix .
docker run --rm --cpus=2 --memory=2g cockroach-130094-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
