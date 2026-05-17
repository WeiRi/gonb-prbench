# gofiber-4062

| Field | Value |
|---|---|
| Project | gofiber |
| Reference | https://github.com/gofiber/gofiber/pull/4062 |
| Category | data_race |
| Oracle | PANIC_RUNTIME |
| Primary diff file | `ctx.go` |
| Base image | `gonb-gofiber-4062-base-v3:latest` (built by gonb-prebuild for this sample) |

**Soft issue — alt-signal oracle**: race condition manifests as `PANIC_RUNTIME` (panic/fatal) rather than race-detector `WARNING: DATA RACE`. Tools that only grep "DATA RACE" will MISS this sample — also grep `panic:` and `fatal error:`.

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t gofiber-4062-bug .
docker run --rm --cpus=2 --memory=2g gofiber-4062-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PANIC_RUNTIME + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t gofiber-4062-fix .
docker run --rm --cpus=2 --memory=2g gofiber-4062-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
