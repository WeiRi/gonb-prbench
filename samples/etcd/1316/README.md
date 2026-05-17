# etcd-1316

| Field | Value |
|---|---|
| Project | etcd |
| Reference | https://github.com/etcd/etcd/pull/1316 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `etcdserver/stats/leader.go` |
| Base image | `gonb-etcd-1316-bug` (built by gonb-prebuild for this sample) |

**Soft issue — dual-test**: `fix.Dockerfile` swaps in a different test file (`verified_test_fixed.go` or `verified_test_inplace.go`) because PR's fix changed an API that the original test referenced. Both bug & fix test the same race, but via different test code.

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t etcd-1316-bug .
docker run --rm --cpus=2 --memory=2g etcd-1316-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t etcd-1316-fix .
docker run --rm --cpus=2 --memory=2g etcd-1316-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
