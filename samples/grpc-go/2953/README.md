# grpc-go-2953

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/2953 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `internal/transport/flowcontrol.go` |
| Base image | `gonb-grpc-go-2953-base-v3:latest` (built by gonb-prebuild for this sample) |

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t grpc-go-2953-bug .
docker run --rm --cpus=2 --memory=2g grpc-go-2953-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t grpc-go-2953-fix .
docker run --rm --cpus=2 --memory=2g grpc-go-2953-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
