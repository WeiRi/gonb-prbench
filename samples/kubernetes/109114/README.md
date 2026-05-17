# kubernetes-109114

| Field | Value |
|---|---|
| Project | kubernetes |
| Reference | https://github.com/kubernetes/kubernetes/pull/109114 |
| Category | data_race |
| Oracle | RACE |
| Primary diff file | `staging/src/k8s.io/client-go/rest/request.go` |
| Base image | `gonb-kubernetes-109114-base:latest` (built by gonb-prebuild for this sample) |

## In-place reproduction

This sample uses the original upstream source at the bug commit (pre-built into the base image — no SSH-clone required at sample-build time).

### Build & run bug

```bash
docker build -f bug.Dockerfile -t kubernetes-109114-bug .
docker run --rm --cpus=2 --memory=2g kubernetes-109114-bug \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: WARNING: DATA RACE + FAIL
```

### Build & run fix

```bash
docker build -f fix.Dockerfile -t kubernetes-109114-fix .
docker run --rm --cpus=2 --memory=2g kubernetes-109114-fix \
  sh -c "go test -race -vet=off -count=10 -timeout=300s ."
# Expected: PASS (PR fix suppresses the race)
```

## Race report

See `race_report_bug.txt` for the captured race detector output from a bug build run.
