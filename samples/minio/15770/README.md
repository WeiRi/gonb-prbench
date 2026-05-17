# minio-15770

| Field | Value |
|---|---|
| Project | minio |
| Reference | https://github.com/minio/minio/pull/15770 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `internal/rest/rpc-stats.go` |
| Base image | `inp-minio-15770` (built by gonb-prebuild for this sample) |

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t minio-15770-bug .
docker run --rm --cpus=2 --memory=2g minio-15770-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t minio-15770-fix .
docker run --rm --cpus=2 --memory=2g minio-15770-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
