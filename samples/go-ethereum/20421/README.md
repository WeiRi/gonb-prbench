# go-ethereum-20421

| Field | Value |
|---|---|
| Project | go-ethereum |
| Reference | https://github.com/ethereum/go-ethereum/pull/20421 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `p2p/enode/iter.go` |
| Base image | `gonb-go-ethereum-20421-bug:latest` (built by gonb-prebuild for this sample) |

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t go-ethereum-20421-bug .
docker run --rm --cpus=2 --memory=2g go-ethereum-20421-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t go-ethereum-20421-fix .
docker run --rm --cpus=2 --memory=2g go-ethereum-20421-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
