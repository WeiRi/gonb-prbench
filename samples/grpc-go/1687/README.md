# grpc-go-1687

| Field | Value |
|---|---|
| Project | grpc-go |
| Reference | https://github.com/grpc/grpc-go/pull/1687 |
| Category | data_race |
| Oracle | PANIC_SEND_CLOSED |
| Primary diff file | `transport/handler_server.go` |
| Base image | `gonb-grpc-go-1687-base-v3:latest` (built by gonb-prebuild for this sample) |

**Soft issue — alt-signal oracle**: race condition manifests as `PANIC_SEND_CLOSED` (panic/fatal) rather than race-detector `WARNING: DATA RACE`. Tools that only grep "DATA RACE" will MISS this sample — also grep `panic:` and `fatal error:`.

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t grpc-go-1687-bug .
docker run --rm --cpus=2 --memory=2g grpc-go-1687-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PANIC_SEND_CLOSED + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t grpc-go-1687-fix .
docker run --rm --cpus=2 --memory=2g grpc-go-1687-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
